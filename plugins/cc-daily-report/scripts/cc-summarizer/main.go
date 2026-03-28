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

type AssistantMsg struct {
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// --- Output structs ---

type SessionSummary struct {
	SessionID         string   `json:"session_id"`
	Project           string   `json:"project"`
	StartTime         string   `json:"start_time"`
	EndTime           string   `json:"end_time"`
	FirstUserMessage  string   `json:"first_user_message"`
	TotalInputTokens  int      `json:"total_input_tokens"`
	TotalOutputTokens int      `json:"total_output_tokens"`
	ToolsUsed         []string `json:"tools_used"`
	TurnCount         int      `json:"turn_count"`
}

type DailyOutput struct {
	Date              string           `json:"date"`
	TotalSessions     int              `json:"total_sessions"`
	TotalInputTokens  int              `json:"total_input_tokens"`
	TotalOutputTokens int              `json:"total_output_tokens"`
	Sessions          []SessionSummary `json:"sessions"`
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

func parseSession(path, targetDate, projectName string) (*SessionSummary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s SessionSummary
	s.Project = projectName
	toolSet := make(map[string]bool)
	firstUser := true

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
		// Filter by date (UTC timestamp prefix)
		if len(entry.Timestamp) >= 10 && entry.Timestamp[:10] != targetDate {
			continue
		}
		if s.SessionID == "" {
			s.SessionID = entry.SessionID
		}
		if s.StartTime == "" {
			s.StartTime = entry.Timestamp
		}
		s.EndTime = entry.Timestamp

		switch entry.Type {
		case "user":
			s.TurnCount++
			if firstUser {
				// Try string content
				var strContent string
				if json.Unmarshal(entry.Message, &struct{ Content *string }{&strContent}) == nil && strContent != "" && strContent != "Warmup" {
					s.FirstUserMessage = truncate(strContent, 200)
					firstUser = false
				} else {
					// Try object with content field
					var obj struct {
						Content interface{} `json:"content"`
					}
					if err := json.Unmarshal(entry.Message, &obj); err == nil {
						switch v := obj.Content.(type) {
						case string:
							if v != "" && v != "Warmup" {
								s.FirstUserMessage = truncate(v, 200)
								firstUser = false
							}
						case []interface{}:
							for _, item := range v {
								if m, ok := item.(map[string]interface{}); ok {
									if m["type"] == "text" {
										if text, ok := m["text"].(string); ok && text != "" {
											s.FirstUserMessage = truncate(text, 200)
											firstUser = false
											break
										}
									}
								}
							}
						}
					}
				}
			}
		case "assistant":
			var msg AssistantMsg
			if err := json.Unmarshal(entry.Message, &msg); err == nil {
				s.TotalInputTokens += msg.Usage.InputTokens
				s.TotalOutputTokens += msg.Usage.OutputTokens
			}
			// Extract tool uses from content array
			var wrapper struct {
				Content []ContentBlock `json:"content"`
			}
			if err := json.Unmarshal(entry.Message, &wrapper); err == nil {
				for _, b := range wrapper.Content {
					if b.Type == "tool_use" && b.Name != "" {
						toolSet[b.Name] = true
					}
				}
			}
		}
	}

	if s.SessionID == "" {
		return nil, nil
	}

	for t := range toolSet {
		s.ToolsUsed = append(s.ToolsUsed, t)
	}
	sort.Strings(s.ToolsUsed)
	if s.ToolsUsed == nil {
		s.ToolsUsed = []string{}
	}
	return &s, nil
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) > max {
		return string(runes[:max]) + "..."
	}
	return s
}

func main() {
	targetDate := time.Now().UTC().Format("2006-01-02")
	if len(os.Args) > 1 {
		targetDate = os.Args[1]
	}

	claudeProjects := filepath.Join(homeDir(), ".claude", "projects")

	// Enumerate project dirs
	projectDirs, err := os.ReadDir(claudeProjects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", claudeProjects, err)
		os.Exit(1)
	}

	type job struct {
		path    string
		project string
	}
	var jobs []job

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
				jobs = append(jobs, job{
					path:    filepath.Join(projPath, e.Name()),
					project: pd.Name(),
				})
			}
		}
	}

	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		sessions []SessionSummary
	)

	for _, j := range jobs {
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			s, err := parseSession(j.path, targetDate, j.project)
			if err != nil || s == nil {
				return
			}
			mu.Lock()
			sessions = append(sessions, *s)
			mu.Unlock()
		}(j)
	}
	wg.Wait()

	// Sort by start time
	sort.Slice(sessions, func(i, k int) bool {
		return sessions[i].StartTime < sessions[k].StartTime
	})

	out := DailyOutput{
		Date:         targetDate,
		TotalSessions: len(sessions),
		Sessions:     sessions,
	}
	if out.Sessions == nil {
		out.Sessions = []SessionSummary{}
	}
	for _, s := range sessions {
		out.TotalInputTokens += s.TotalInputTokens
		out.TotalOutputTokens += s.TotalOutputTokens
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}
