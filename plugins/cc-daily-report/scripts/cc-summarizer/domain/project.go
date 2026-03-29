package domain

import (
	"regexp"
	"strings"
)

var notionSectionRe = regexp.MustCompile(`(?m)##\s*📚\s*Notion\s*\n+\s*(https://www\.notion\.so/\S+)`)
var notionFallbackRe = regexp.MustCompile(`https://www\.notion\.so/\S+`)

// PRInfo は GitHub PR の情報を表す。
type PRInfo struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	NotionURL string `json:"notion_url,omitempty"`
}

// NewPRInfo は PR の各フィールドと本文から PRInfo を生成する。
// Notion URL の抽出は内部で行う。
func NewPRInfo(number int, title, url, body string) *PRInfo {
	return &PRInfo{
		Number:    number,
		Title:     title,
		URL:       url,
		NotionURL: extractNotionURL(body),
	}
}

// extractNotionURL は PR 本文から Notion URL を抽出する。
func extractNotionURL(body string) string {
	if m := notionSectionRe.FindStringSubmatch(body); len(m) >= 2 {
		return strings.TrimRight(m[1], ")")
	}
	return notionFallbackRe.FindString(body)
}

// ---

// PRTimeSummary は PR ごとの作業時間集計を表す。
type PRTimeSummary struct {
	PRNumber   int     `json:"pr_number"`
	PRTitle    string  `json:"pr_title"`
	PRURL      string  `json:"pr_url"`
	NotionURL  string  `json:"notion_url,omitempty"`
	TotalHours float64 `json:"total_hours"`
}

// ---

// RepositorySummary はリポジトリ（プロジェクト）単位の集計を表す。
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
	PRTimeSummary     []PRTimeSummary  `json:"pr_time_summary,omitempty"`
	Sessions          []SessionSummary `json:"sessions"`
}

// ---

// DailyOutput は日次集計の最終出力構造体。
type DailyOutput struct {
	Date              string              `json:"date"`
	TotalRepositories int                 `json:"total_repositories"`
	TotalSessions     int                 `json:"total_sessions"`
	TotalInputTokens  int                 `json:"total_input_tokens"`
	TotalOutputTokens int                 `json:"total_output_tokens"`
	Repositories      []RepositorySummary `json:"repositories"`
}
