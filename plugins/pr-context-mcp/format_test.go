package main

import (
	"strings"
	"testing"
)

func TestFormatOutput_NoPRs(t *testing.T) {
	result := formatOutput([]PRDetails{})
	if !strings.Contains(result, "見つかりませんでした") {
		t.Errorf("PRなしのメッセージが含まれるべきです:\n%s", result)
	}
}

func TestFormatOutput_SinglePR(t *testing.T) {
	details := []PRDetails{
		{
			PR: PullRequest{
				Number:   123,
				Title:    "テスト修正",
				Body:     "バグを修正しました",
				State:    "closed",
				MergedAt: "2024-03-01T12:00:00Z",
			},
			ReviewComments: []Comment{
				{Body: "LGTMです", User: User{Login: "reviewer"}},
			},
			IssueComments: []Comment{},
		},
	}

	result := formatOutput(details)

	if !strings.Contains(result, "PR #123") {
		t.Errorf("PR番号が含まれるべきです:\n%s", result)
	}
	if !strings.Contains(result, "テスト修正") {
		t.Errorf("PRタイトルが含まれるべきです:\n%s", result)
	}
	if !strings.Contains(result, "バグを修正しました") {
		t.Errorf("ディスクリプションが含まれるべきです:\n%s", result)
	}
	if !strings.Contains(result, "2024-03-01") {
		t.Errorf("マージ日（日付部分）が含まれるべきです:\n%s", result)
	}
	if !strings.Contains(result, "@reviewer") {
		t.Errorf("レビュアー名が含まれるべきです:\n%s", result)
	}
	if !strings.Contains(result, "LGTMです") {
		t.Errorf("レビューコメントが含まれるべきです:\n%s", result)
	}
}

func TestFormatOutput_MultiplePRs(t *testing.T) {
	details := []PRDetails{
		{
			PR: PullRequest{Number: 1, Title: "1つ目のPR", State: "closed"},
		},
		{
			PR: PullRequest{Number: 2, Title: "2つ目のPR", State: "open"},
		},
	}

	result := formatOutput(details)

	if !strings.Contains(result, "PR #1") {
		t.Errorf("1つ目のPR番号が含まれるべきです:\n%s", result)
	}
	if !strings.Contains(result, "PR #2") {
		t.Errorf("2つ目のPR番号が含まれるべきです:\n%s", result)
	}
	// 区切り線が含まれていること
	if !strings.Contains(result, "---") {
		t.Errorf("複数PRの区切りが含まれるべきです:\n%s", result)
	}
}

func TestFormatOutput_NoComments(t *testing.T) {
	details := []PRDetails{
		{
			PR: PullRequest{
				Number: 5,
				Title:  "コメントなしPR",
				Body:   "変更の説明",
				State:  "merged",
			},
			ReviewComments: []Comment{},
			IssueComments:  []Comment{},
		},
	}

	result := formatOutput(details)

	if strings.Contains(result, "レビューコメント") {
		t.Errorf("コメントがない場合、レビューコメントセクションは表示されるべきでありません:\n%s", result)
	}
	if strings.Contains(result, "一般コメント") {
		t.Errorf("コメントがない場合、一般コメントセクションは表示されるべきでありません:\n%s", result)
	}
}

func TestFormatOutput_EmptyBody(t *testing.T) {
	details := []PRDetails{
		{
			PR: PullRequest{
				Number: 7,
				Title:  "本文なしPR",
				Body:   "",
				State:  "closed",
			},
		},
	}

	result := formatOutput(details)

	if strings.Contains(result, "ディスクリプション") {
		t.Errorf("本文が空の場合、ディスクリプションセクションは表示されるべきでありません:\n%s", result)
	}
}

func TestFormatOutput_MergedAtTruncated(t *testing.T) {
	details := []PRDetails{
		{
			PR: PullRequest{
				Number:   10,
				Title:    "マージ済みPR",
				State:    "closed",
				MergedAt: "2024-06-15T08:30:00Z",
			},
		},
	}

	result := formatOutput(details)

	// タイムスタンプ全体ではなく日付部分のみ表示されていること
	if strings.Contains(result, "T08:30:00Z") {
		t.Errorf("時刻部分は表示されるべきでありません:\n%s", result)
	}
	if !strings.Contains(result, "2024-06-15") {
		t.Errorf("日付部分が含まれるべきです:\n%s", result)
	}
}
