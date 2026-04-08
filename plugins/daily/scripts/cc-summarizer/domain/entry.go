package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RawEntry は JSONL の 1 行分の生データを表す。
type RawEntry struct {
	Type      string          `json:"type"`
	SessionID string          `json:"sessionId"`
	Timestamp string          `json:"timestamp"`
	CWD       string          `json:"cwd"`
	Message   json.RawMessage `json:"message"`
}

// IsForDate は、このエントリのタイムスタンプが指定日（YYYY-MM-DD）に一致するか返す。
func (e *RawEntry) IsForDate(date string) bool {
	return len(e.Timestamp) >= 10 && e.Timestamp[:10] == date
}

// UserText は user エントリの Message フィールドからプレーンテキストを抽出する。
func (e *RawEntry) UserText() string {
	var s string
	if json.Unmarshal(e.Message, &s) == nil {
		return summarizeSkillContent(s)
	}
	var obj struct {
		Content interface{} `json:"content"`
	}
	if json.Unmarshal(e.Message, &obj) != nil {
		return ""
	}
	switch v := obj.Content.(type) {
	case string:
		return summarizeSkillContent(v)
	case []interface{}:
		var parts []string
		for _, item := range v {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if m["type"] == "text" {
				if text, ok := m["text"].(string); ok && text != "" {
					parts = append(parts, summarizeSkillContent(text))
				}
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

// summarizeSkillContent は "Base directory for this skill" で始まる内容を短縮する。
func summarizeSkillContent(text string) string {
	if strings.HasPrefix(text, "Base directory for this skill:") {
		lines := strings.SplitN(text, "\n", 2)
		skillPath := strings.TrimPrefix(lines[0], "Base directory for this skill: ")
		skillName := filepath.Base(skillPath)
		return fmt.Sprintf("[skill: %sを使用]", skillName)
	}
	return text
}

// IsSystemMessage は、このユーザーエントリが自動生成メッセージ（スラッシュコマンド出力・Warmup 等）
// かどうかを返す。true の場合は会話履歴に含めない。
func (e *RawEntry) IsSystemMessage() bool {
	text := e.UserText()
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

// projectNameFromCWD は作業ディレクトリのパスからプロジェクト名を返す。
// パッケージ内で共有するプライベートヘルパー。
func projectNameFromCWD(cwd string) string {
	if cwd == "" {
		return "unknown"
	}
	// cwdから親ディレクトリを辿って .git を探す
	dir := cwd
	for {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			if info.IsDir() {
				// 通常のgitリポジトリ → そのディレクトリ名を返す
				return filepath.Base(dir)
			}
			// git worktree: .git はファイルで "gitdir: /path/to/.git/worktrees/xxx" を含む
			data, err := os.ReadFile(gitPath)
			if err != nil {
				return filepath.Base(dir)
			}
			line := strings.TrimSpace(string(data))
			if !strings.HasPrefix(line, "gitdir: ") {
				return filepath.Base(dir)
			}
			gitdir := strings.TrimPrefix(line, "gitdir: ")
			if !filepath.IsAbs(gitdir) {
				gitdir = filepath.Join(dir, gitdir)
			}
			// gitdir は "/path/to/repo/.git/worktrees/xxx" の形式
			// ".git/worktrees" を含むなら、その2つ上が親リポジトリの .git
			cleaned := filepath.Clean(gitdir)
			if i := strings.Index(cleaned, string(filepath.Separator)+".git"+string(filepath.Separator)+"worktrees"+string(filepath.Separator)); i >= 0 {
				return filepath.Base(cleaned[:i])
			}
			return filepath.Base(dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return filepath.Base(cwd)
}

// ---

// ContentBlock は assistant メッセージ内のコンテンツブロックを表す。
type ContentBlock struct {
	Type  string          `json:"type"`
	Name  string          `json:"name,omitempty"`
	Text  string          `json:"text,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

// Files は tool_use ブロックの Input からファイルパスを抽出して返す。
func (c *ContentBlock) Files() []string {
	var m map[string]interface{}
	if json.Unmarshal(c.Input, &m) != nil {
		return nil
	}
	var paths []string
	switch c.Name {
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

// ---

// AssistantMsg は assistant エントリの Message フィールドをパースした構造体。
type AssistantMsg struct {
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Content []ContentBlock `json:"content"`
}
