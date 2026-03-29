package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// --- Raw JSONL structs ---

type RawEntry struct {
	Type      string          `json:"type"`
	SessionID string          `json:"sessionId"`
	Timestamp string          `json:"timestamp"`
	CWD       string          `json:"cwd"`
	Message   json.RawMessage `json:"message"`
}

type ContentBlock struct {
	Type  string          `json:"type"`
	Name  string          `json:"name,omitempty"`
	Text  string          `json:"text,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

type AssistantMsg struct {
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Content []ContentBlock `json:"content"`
}

// --- Output structs ---

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SessionSummary struct {
	SessionID         string         `json:"session_id"`
	Project           string         `json:"project"`
	CWD               string         `json:"-"`
	GitBranch         string         `json:"git_branch,omitempty"`
	StartTime         string         `json:"start_time"`
	EndTime           string         `json:"end_time"`
	Conversation      []Message      `json:"conversation"`
	TotalInputTokens  int            `json:"total_input_tokens"`
	TotalOutputTokens int            `json:"total_output_tokens"`
	ToolsUsed         map[string]int `json:"tools_used"`
	FilesAccessed     []string       `json:"files_accessed"`
}

type RepositorySummary struct {
	Project           string           `json:"project"`
	TotalSessions     int              `json:"total_sessions"`
	TotalInputTokens  int              `json:"total_input_tokens"`
	TotalOutputTokens int              `json:"total_output_tokens"`
	StartTime         string           `json:"start_time"`
	EndTime           string           `json:"end_time"`
	GitBranches       []string         `json:"git_branches"`
	ToolsUsed         map[string]int   `json:"tools_used"`
	FilesAccessed     []string         `json:"files_accessed"`
	Sessions          []SessionSummary `json:"sessions"`
}

type DailyOutput struct {
	Date              string              `json:"date"`
	TotalRepositories int                 `json:"total_repositories"`
	TotalSessions     int                 `json:"total_sessions"`
	TotalInputTokens  int                 `json:"total_input_tokens"`
	TotalOutputTokens int                 `json:"total_output_tokens"`
	Repositories      []RepositorySummary `json:"repositories"`
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("HOME")
	}
	return home
}

func projectNameFromCWD(cwd string) string {
	if cwd == "" {
		return "unknown"
	}
	return filepath.Base(cwd)
}

func gitBranchFromCWD(cwd string) string {
	if cwd == "" {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(cwd, ".git", "HEAD"))
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	if strings.HasPrefix(line, "ref: refs/heads/") {
		return strings.TrimPrefix(line, "ref: refs/heads/")
	}
	if len(line) >= 7 {
		return line[:7]
	}
	return line
}

// isSystemMessage returns true for auto-generated messages that should be excluded
// from the conversation history (slash command output, warmup, etc.).
func isSystemMessage(text string) bool {
	if text == "" || text == "Warmup" {
		return true
	}
	for _, prefix := range []string{
		"<local-command-caveat>",
		"<command-name>",
		"<local-command-stdout>",
		"<command-message>",
	} {
		if strings.HasPrefix(text, prefix) {
			return true
		}
	}
	return false
}

func extractUserText(msg json.RawMessage) string {
	var s string
	if json.Unmarshal(msg, &s) == nil {
		return s
	}
	var obj struct {
		Content interface{} `json:"content"`
	}
	if json.Unmarshal(msg, &obj) != nil {
		return ""
	}
	switch v := obj.Content.(type) {
	case string:
		return v
	case []interface{}:
		var parts []string
		for _, item := range v {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if m["type"] == "text" {
				if text, ok := m["text"].(string); ok && text != "" {
					parts = append(parts, text)
				}
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

func extractFilesFromInput(toolName string, input json.RawMessage) []string {
	var m map[string]interface{}
	if json.Unmarshal(input, &m) != nil {
		return nil
	}
	var paths []string
	switch toolName {
	case "Read", "Edit", "Write":
		if v, ok := m["file_path"].(string); ok && v != "" {
			paths = append(paths, v)
		}
	case "Grep", "Glob":
		if v, ok := m["path"].(string); ok && v != "" {
			paths = append(paths, v)
		}
	}
	return paths
}

func parseSession(path, targetDate string) (*SessionSummary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s SessionSummary
	s.ToolsUsed = make(map[string]int)
	fileSet := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 10*1024*1024), 10*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry RawEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if entry.Type == "file-history-snapshot" || entry.Timestamp == "" {
			continue
		}
		if len(entry.Timestamp) >= 10 && entry.Timestamp[:10] != targetDate {
			continue
		}

		if s.SessionID == "" {
			s.SessionID = entry.SessionID
		}
		if s.CWD == "" && entry.CWD != "" {
			s.CWD = entry.CWD
		}
		if s.StartTime == "" {
			s.StartTime = entry.Timestamp
		}
		s.EndTime = entry.Timestamp

		switch entry.Type {
		case "user":
			text := extractUserText(entry.Message)
			if !isSystemMessage(text) {
				s.Conversation = append(s.Conversation, Message{Role: "user", Content: text})
			}

		case "assistant":
			var msg AssistantMsg
			if err := json.Unmarshal(entry.Message, &msg); err != nil {
				continue
			}
			s.TotalInputTokens += msg.Usage.InputTokens
			s.TotalOutputTokens += msg.Usage.OutputTokens

			var textParts []string
			for _, b := range msg.Content {
				switch b.Type {
				case "text":
					if b.Text != "" {
						textParts = append(textParts, b.Text)
					}
				case "tool_use":
					if b.Name != "" {
						s.ToolsUsed[b.Name]++
						for _, fp := range extractFilesFromInput(b.Name, b.Input) {
							fileSet[fp] = true
						}
					}
				}
			}
			if len(textParts) > 0 {
				s.Conversation = append(s.Conversation, Message{
					Role:    "assistant",
					Content: strings.Join(textParts, "\n"),
				})
			}
		}
	}

	if s.SessionID == "" || len(s.Conversation) == 0 {
		return nil, nil
	}

	s.Project = projectNameFromCWD(s.CWD)
	s.GitBranch = gitBranchFromCWD(s.CWD)

	for fp := range fileSet {
		s.FilesAccessed = append(s.FilesAccessed, fp)
	}
	sort.Strings(s.FilesAccessed)
	if s.FilesAccessed == nil {
		s.FilesAccessed = []string{}
	}

	return &s, nil
}

func main() {
	targetDate := time.Now().UTC().Format("2006-01-02")
	if len(os.Args) > 1 {
		targetDate = os.Args[1]
	}

	claudeProjects := filepath.Join(homeDir(), ".claude", "projects")

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
		sessions []SessionSummary
	)

	for _, p := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			s, err := parseSession(p, targetDate)
			if err != nil || s == nil {
				return
			}
			mu.Lock()
			sessions = append(sessions, *s)
			mu.Unlock()
		}(p)
	}
	wg.Wait()

	sort.Slice(sessions, func(i, k int) bool {
		return sessions[i].StartTime < sessions[k].StartTime
	})

	// Group sessions by project name
	repoMap := make(map[string]*RepositorySummary)
	var repoOrder []string

	for _, s := range sessions {
		key := s.Project
		if _, exists := repoMap[key]; !exists {
			repoMap[key] = &RepositorySummary{
				Project:   key,
				ToolsUsed: make(map[string]int),
			}
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
		if s.GitBranch != "" {
			found := false
			for _, b := range r.GitBranches {
				if b == s.GitBranch {
					found = true
					break
				}
			}
			if !found {
				r.GitBranches = append(r.GitBranches, s.GitBranch)
			}
		}
		for tool, count := range s.ToolsUsed {
			r.ToolsUsed[tool] += count
		}
		fileSet := make(map[string]bool)
		for _, fp := range r.FilesAccessed {
			fileSet[fp] = true
		}
		for _, fp := range s.FilesAccessed {
			fileSet[fp] = true
		}
		r.FilesAccessed = nil
		for fp := range fileSet {
			r.FilesAccessed = append(r.FilesAccessed, fp)
		}
		sort.Strings(r.FilesAccessed)
		r.Sessions = append(r.Sessions, s)
	}

	var repos []RepositorySummary
	for _, key := range repoOrder {
		r := repoMap[key]
		if r.GitBranches == nil {
			r.GitBranches = []string{}
		}
		if r.FilesAccessed == nil {
			r.FilesAccessed = []string{}
		}
		repos = append(repos, *r)
	}
	if repos == nil {
		repos = []RepositorySummary{}
	}

	out := DailyOutput{
		Date:              targetDate,
		TotalRepositories: len(repos),
		TotalSessions:     len(sessions),
		Repositories:      repos,
	}
	for _, r := range repos {
		out.TotalInputTokens += r.TotalInputTokens
		out.TotalOutputTokens += r.TotalOutputTokens
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}
