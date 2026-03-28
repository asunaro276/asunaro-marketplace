package main

import (
	"errors"
	"strings"
	"testing"
)

// --- モック実装 ---

type mockGitClient struct {
	sha       string
	shaErr    error
	owner     string
	repo      string
	remoteErr error
}

func (m *mockGitClient) BlameCommit(cwd, filePath string, line int) (string, error) {
	return m.sha, m.shaErr
}

func (m *mockGitClient) RemoteRepo(cwd string) (string, string, error) {
	return m.owner, m.repo, m.remoteErr
}

type mockGitHubClient struct {
	prs        []PullRequest
	prsErr     error
	details    PRDetails
	detailsErr error
}

func (m *mockGitHubClient) GetPRsForCommit(owner, repo, sha string) ([]PullRequest, error) {
	return m.prs, m.prsErr
}

func (m *mockGitHubClient) GetPRDetails(owner, repo string, number int) (PRDetails, error) {
	return m.details, m.detailsErr
}

// --- テスト ---

func TestGetPRContext_Success(t *testing.T) {
	svc := newService(
		&mockGitClient{sha: "abc123def456abc123def456abc123def456abc1", owner: "owner", repo: "repo"},
		&mockGitHubClient{
			prs: []PullRequest{{Number: 42, Title: "バグ修正"}},
			details: PRDetails{
				PR: PullRequest{
					Number:   42,
					Title:    "バグ修正",
					Body:     "nilポインタ参照を修正しました",
					State:    "closed",
					MergedAt: "2024-01-15T10:00:00Z",
				},
				ReviewComments: []Comment{
					{Body: "良い修正です", User: User{Login: "reviewer", Type: "User"}},
				},
				IssueComments: []Comment{},
			},
		},
	)

	result, err := svc.GetPRContext("main.go", 10, "", "/tmp/repo")
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if !strings.Contains(result, "PR #42") {
		t.Errorf("PR番号が含まれていません:\n%s", result)
	}
	if !strings.Contains(result, "バグ修正") {
		t.Errorf("PRタイトルが含まれていません:\n%s", result)
	}
	if !strings.Contains(result, "nilポインタ参照を修正しました") {
		t.Errorf("PRディスクリプションが含まれていません:\n%s", result)
	}
	if !strings.Contains(result, "良い修正です") {
		t.Errorf("レビューコメントが含まれていません:\n%s", result)
	}
}

func TestGetPRContext_NoPRs(t *testing.T) {
	svc := newService(
		&mockGitClient{sha: "abc123def456abc123def456abc123def456abc1", owner: "owner", repo: "repo"},
		&mockGitHubClient{prs: []PullRequest{}},
	)

	result, err := svc.GetPRContext("main.go", 1, "", "/tmp/repo")
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if !strings.Contains(result, "直接コミット") {
		t.Errorf("直接コミットのメッセージが含まれていません:\n%s", result)
	}
}

func TestGetPRContext_UncommittedLine(t *testing.T) {
	svc := newService(
		&mockGitClient{sha: "0000000000000000000000000000000000000000", owner: "owner", repo: "repo"},
		&mockGitHubClient{},
	)

	result, err := svc.GetPRContext("main.go", 1, "", "/tmp/repo")
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if !strings.Contains(result, "コミットされていない") {
		t.Errorf("未コミットのメッセージが含まれていません:\n%s", result)
	}
}

func TestGetPRContext_GitBlameError(t *testing.T) {
	svc := newService(
		&mockGitClient{shaErr: errors.New("not a git repository")},
		&mockGitHubClient{},
	)

	_, err := svc.GetPRContext("main.go", 1, "", "/tmp/repo")
	if err == nil {
		t.Error("エラーが返されるべきです")
	}
	if !strings.Contains(err.Error(), "git blame") {
		t.Errorf("エラーメッセージに 'git blame' が含まれていません: %v", err)
	}
}

func TestGetPRContext_ExplicitRepo(t *testing.T) {
	var capturedOwner, capturedRepo string
	gh := &mockGitHubClient{
		prs: []PullRequest{},
	}
	// RemoteRepo が呼ばれないことを確認するため、呼ばれたらエラーにするモック
	git := &mockGitClient{
		sha:       "abc123def456abc123def456abc123def456abc1",
		remoteErr: errors.New("RemoteRepo should not be called"),
	}
	_ = capturedOwner
	_ = capturedRepo

	svc := newService(git, gh)
	result, err := svc.GetPRContext("main.go", 1, "myorg/myrepo", "/tmp/repo")
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	// PR がない場合のメッセージが返れば、repo の解析が正常に行われた証拠
	if !strings.Contains(result, "直接コミット") {
		t.Errorf("想定外の結果:\n%s", result)
	}
}

func TestGetPRContext_InvalidRepoFormat(t *testing.T) {
	svc := newService(
		&mockGitClient{sha: "abc123def456abc123def456abc123def456abc1"},
		&mockGitHubClient{},
	)

	_, err := svc.GetPRContext("main.go", 1, "invalid-format", "/tmp/repo")
	if err == nil {
		t.Error("不正な repo フォーマットでエラーが返されるべきです")
	}
}

func TestGetPRContext_BotCommentsFiltered(t *testing.T) {
	svc := newService(
		&mockGitClient{sha: "abc123def456abc123def456abc123def456abc1", owner: "owner", repo: "repo"},
		&mockGitHubClient{
			prs: []PullRequest{{Number: 1, Title: "修正"}},
			details: PRDetails{
				PR: PullRequest{Number: 1, Title: "修正", State: "closed"},
				ReviewComments: []Comment{
					{Body: "人間のコメント", User: User{Login: "developer", Type: "User"}},
					{Body: "カバレッジレポート", User: User{Login: "codecov", Type: "Bot"}},
				},
				IssueComments: []Comment{
					{Body: "CI通過しました", User: User{Login: "github-actions[bot]", Type: "Bot"}},
				},
			},
		},
	)

	result, err := svc.GetPRContext("main.go", 1, "", "/tmp/repo")
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if !strings.Contains(result, "人間のコメント") {
		t.Errorf("人間のコメントが含まれるべきです:\n%s", result)
	}
	if strings.Contains(result, "カバレッジレポート") {
		t.Errorf("ボットのコメントは除外されるべきです:\n%s", result)
	}
	if strings.Contains(result, "CI通過しました") {
		t.Errorf("ボットのコメントは除外されるべきです:\n%s", result)
	}
}
