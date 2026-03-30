package domain

import (
	"encoding/json"
	"testing"
)

func TestSummarizeSkillContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "skill content is summarized",
			input:    "Base directory for this skill: /home/user/.claude/plugins/git-operation/skills/commit\n\n# Git Commit\n\ngit-operations-expertサブエージェントを使用して...",
			expected: "[skill: commitを使用]",
		},
		{
			name:     "daily-report skill is summarized",
			input:    "Base directory for this skill: /home/user/.claude/plugins/cc-daily-report/skills/daily-report\n\n# Claude Code 日報生成スキル\n\n## 概要...",
			expected: "[skill: daily-reportを使用]",
		},
		{
			name:     "non-skill content is unchanged",
			input:    "This is a normal user message",
			expected: "This is a normal user message",
		},
		{
			name:     "empty string is unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "create-pr skill is summarized",
			input:    "Base directory for this skill: /home/user/.claude/plugins/git-operation/skills/create-pr\n\n# Git Create PR...",
			expected: "[skill: create-prを使用]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := summarizeSkillContent(tc.input)
			if result != tc.expected {
				t.Errorf("got %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestUserTextWithSkillContent(t *testing.T) {
	skillContent := "Base directory for this skill: /home/user/.claude/plugins/git-operation/skills/commit\n\n# Git Commit\n\nLong content here..."
	msg, _ := json.Marshal(skillContent)
	entry := RawEntry{
		Type:      "user",
		SessionID: "test-session",
		Timestamp: "2024-06-15T10:00:00Z",
		CWD:       "/home/user/myapp",
		Message:   msg,
	}

	result := entry.UserText()
	expected := "[skill: commitを使用]"
	if result != expected {
		t.Errorf("UserText() got %q, want %q", result, expected)
	}
}

func TestUserTextWithArrayContentContainingSkill(t *testing.T) {
	skillContent := "Base directory for this skill: /home/user/.claude/plugins/cc-daily-report/skills/daily-report\n\n# Daily Report..."
	content := []interface{}{
		map[string]interface{}{
			"type": "text",
			"text": skillContent,
		},
	}
	msg, _ := json.Marshal(map[string]interface{}{
		"role":    "user",
		"content": content,
	})
	entry := RawEntry{
		Type:      "user",
		SessionID: "test-session",
		Timestamp: "2024-06-15T10:00:00Z",
		CWD:       "/home/user/myapp",
		Message:   msg,
	}

	result := entry.UserText()
	expected := "[skill: daily-reportを使用]"
	if result != expected {
		t.Errorf("UserText() got %q, want %q", result, expected)
	}
}

func TestUserTextWithMixedContent(t *testing.T) {
	skillContent := "Base directory for this skill: /home/user/.claude/plugins/git-operation/skills/commit\n\n# Git Commit..."
	normalContent := "This is a normal message"
	content := []interface{}{
		map[string]interface{}{
			"type": "text",
			"text": skillContent,
		},
		map[string]interface{}{
			"type": "text",
			"text": normalContent,
		},
	}
	msg, _ := json.Marshal(map[string]interface{}{
		"role":    "user",
		"content": content,
	})
	entry := RawEntry{
		Type:      "user",
		SessionID: "test-session",
		Timestamp: "2024-06-15T10:00:00Z",
		CWD:       "/home/user/myapp",
		Message:   msg,
	}

	result := entry.UserText()
	expected := "[skill: commitを使用]\nThis is a normal message"
	if result != expected {
		t.Errorf("UserText() got %q, want %q", result, expected)
	}
}
