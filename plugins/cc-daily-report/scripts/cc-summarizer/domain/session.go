package domain

import (
	"encoding/json"
	"sort"
	"strings"
)

// Message は会話の 1 ターンを表す出力用の構造体。
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SessionSummary は 1 セッション分の解析結果を表す。
type SessionSummary struct {
	SessionID         string         `json:"session_id"`
	Project           string         `json:"project"`
	CWD               string         `json:"-"`
	GitBranch         string         `json:"git_branch,omitempty"`
	StartTime         string         `json:"start_time"`
	EndTime           string         `json:"end_time"`
	DurationMinutes   float64        `json:"duration_minutes"`
	PRInfo            *PRInfo        `json:"pr_info,omitempty"`
	Conversation      []Message      `json:"conversation"`
	TotalInputTokens  int            `json:"total_input_tokens"`
	TotalOutputTokens int            `json:"total_output_tokens"`
	ToolsUsed         map[string]int `json:"tools_used"`
	FilesAccessed     []string       `json:"files_accessed"`
}

// NewSession はフィルタ済みの RawEntry スライスから SessionSummary を構築する。
// エントリは呼び出し前に日付・型でフィルタ済みであること。
// セッションが成立しない場合は nil を返す。
func NewSession(entries []RawEntry) *SessionSummary {
	var s SessionSummary
	s.ToolsUsed = make(map[string]int)
	fileSet := make(map[string]bool)

	for _, entry := range entries {
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
			if !entry.IsSystemMessage() {
				s.Conversation = append(s.Conversation, Message{Role: "user", Content: entry.UserText()})
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
						for _, fp := range b.Files() {
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
		return nil
	}

	s.Project = projectNameFromCWD(s.CWD)
	s.DurationMinutes = SessionDuration(s.StartTime, s.EndTime)

	for fp := range fileSet {
		s.FilesAccessed = append(s.FilesAccessed, fp)
	}
	sort.Strings(s.FilesAccessed)
	if s.FilesAccessed == nil {
		s.FilesAccessed = []string{}
	}

	return &s
}
