package main

import (
	"cc-summarizer/domain"
	"cc-summarizer/infra"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

func main() {
	targetDate := time.Now().UTC().Format("2006-01-02")
	if len(os.Args) > 1 {
		targetDate = os.Args[1]
	}

	claudeProjects := filepath.Join(infra.HomeDir(), ".claude", "projects")

	projectDirs, err := os.ReadDir(claudeProjects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", claudeProjects, err)
		os.Exit(1)
	}

	var paths []string
	for _, pd := range projectDirs {
		if !pd.IsDir() {
			continue
		}
		projPath := filepath.Join(claudeProjects, pd.Name())
		entries, _ := os.ReadDir(projPath)
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
				continue
			}
			info, _ := e.Info()
			if info != nil && info.ModTime().UTC().Format("2006-01-02") == targetDate {
				paths = append(paths, filepath.Join(projPath, e.Name()))
			}
		}
	}

	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		sessions []domain.SessionSummary
	)

	for _, p := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			s, err := infra.ReadSession(p, targetDate)
			if err != nil || s == nil {
				return
			}
			mu.Lock()
			sessions = append(sessions, *s)
			mu.Unlock()
		}(p)
	}
	wg.Wait()

	// Fetch PR info for all unique cwd+branch combos in parallel
	type cwdBranch struct{ cwd, branch string }
	keySet := make(map[cwdBranch]bool)
	for i := range sessions {
		s := &sessions[i]
		if s.CWD != "" && s.GitBranch != "" {
			keySet[cwdBranch{s.CWD, s.GitBranch}] = true
		}
	}
	var prMu sync.Mutex
	var prWg sync.WaitGroup
	prCache := make(map[string]*domain.PRInfo)
	for k := range keySet {
		prWg.Add(1)
		go func(k cwdBranch) {
			defer prWg.Done()
			info := infra.FetchPRInfo(k.cwd, k.branch)
			prMu.Lock()
			prCache[k.cwd+"|"+k.branch] = info
			prMu.Unlock()
		}(k)
	}
	prWg.Wait()
	for i := range sessions {
		s := &sessions[i]
		if s.CWD != "" && s.GitBranch != "" {
			s.PRInfo = prCache[s.CWD+"|"+s.GitBranch]
		}
		s.ClassifyWorkImportance()
	}

	sort.Slice(sessions, func(i, k int) bool {
		return sessions[i].StartTime < sessions[k].StartTime
	})

	// Group sessions by project name
	repoMap := make(map[string]*domain.RepositorySummary)
	repoFileSets := make(map[string]map[string]bool)
	repoBranchSets := make(map[string]map[string]bool)
	repoPRTimeMaps := make(map[string]map[int]*domain.PRTimeSummary)
	repoPROrders := make(map[string][]int)
	var repoOrder []string

	for _, s := range sessions {
		key := s.Project
		if _, exists := repoMap[key]; !exists {
			repoMap[key] = &domain.RepositorySummary{
				Project:   key,
				ToolsUsed: make(map[string]int),
			}
			repoFileSets[key] = make(map[string]bool)
			repoBranchSets[key] = make(map[string]bool)
			repoPRTimeMaps[key] = make(map[int]*domain.PRTimeSummary)
			repoOrder = append(repoOrder, key)
		}
		r := repoMap[key]
		r.TotalSessions++
		r.TotalInputTokens += s.TotalInputTokens
		r.TotalOutputTokens += s.TotalOutputTokens
		if r.StartTime == "" || s.StartTime < r.StartTime {
			r.StartTime = s.StartTime
		}
		if s.EndTime > r.EndTime {
			r.EndTime = s.EndTime
		}
		if s.GitBranch != "" && !repoBranchSets[key][s.GitBranch] {
			repoBranchSets[key][s.GitBranch] = true
			r.GitBranches = append(r.GitBranches, s.GitBranch)
		}
		for tool, count := range s.ToolsUsed {
			r.ToolsUsed[tool] += count
		}
		for _, fp := range s.FilesAccessed {
			repoFileSets[key][fp] = true
		}
		if s.PRInfo != nil {
			pr := s.PRInfo
			if _, exists := repoPRTimeMaps[key][pr.Number]; !exists {
				repoPRTimeMaps[key][pr.Number] = &domain.PRTimeSummary{
					PRNumber:  pr.Number,
					PRTitle:   pr.Title,
					PRURL:     pr.URL,
					NotionURL: pr.NotionURL,
				}
				repoPROrders[key] = append(repoPROrders[key], pr.Number)
			}
			repoPRTimeMaps[key][pr.Number].TotalHours += s.DurationMinutes / 60.0
		}
		r.Sessions = append(r.Sessions, s)
	}

	var repos []domain.RepositorySummary
	for _, key := range repoOrder {
		r := repoMap[key]
		for fp := range repoFileSets[key] {
			r.FilesAccessed = append(r.FilesAccessed, fp)
		}
		sort.Strings(r.FilesAccessed)
		for _, num := range repoPROrders[key] {
			r.PRTimeSummary = append(r.PRTimeSummary, *repoPRTimeMaps[key][num])
		}
		if r.GitBranches == nil {
			r.GitBranches = []string{}
		}
		if r.FilesAccessed == nil {
			r.FilesAccessed = []string{}
		}
		repos = append(repos, *r)
	}
	if repos == nil {
		repos = []domain.RepositorySummary{}
	}

	out := domain.DailyOutput{
		Date:                     targetDate,
		TotalRepositories:        len(repos),
		TotalSessions:            len(sessions),
		Repositories:             repos,
		RecommendedQuestionCount: domain.CalcRecommendedQuestionCount(len(sessions)),
	}
	for _, r := range repos {
		out.TotalInputTokens += r.TotalInputTokens
		out.TotalOutputTokens += r.TotalOutputTokens
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}
