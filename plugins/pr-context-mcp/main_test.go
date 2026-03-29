package main

import (
	"errors"
	"strings"
	"testing"
)

// ============================================================
// git.go のテスト
// ============================================================

func TestExtractSHAFromBlame_Success(t *testing.T) {
	// git blame --porcelain の典型的な出力
	output := "abc123def456abc123def456abc123def456abc1 10 10 1\nauthor Alice\n"
	sha, err := extractSHAFromBlame(output)
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if sha != "abc123def456abc123def456abc123def456abc1" {
		t.Errorf("期待する SHA と異なります: %s", sha)
	}
}

func TestExtractSHAFromBlame_InvalidOutput(t *testing.T) {
	_, err := extractSHAFromBlame("")
	if err == nil {
		t.Error("空の出力でエラーが返されるべきです")
	}
}

func TestParseRemoteURL_HTTPS(t *testing.T) {
	cases := []struct {
		url         string
		wantOwner   string
		wantRepo    string
	}{
		{"https://github.com/owner/repo.git", "owner", "repo"},
		{"https://github.com/owner/repo", "owner", "repo"},
		{"http://github.com/myorg/myproject.git", "myorg", "myproject"},
	}
	for _, tc := range cases {
		owner, repo, err := parseRemoteURL(tc.url)
		if err != nil {
			t.Errorf("URL %q でエラーが発生しました: %v", tc.url, err)
			continue
		}
		if owner != tc.wantOwner || repo != tc.wantRepo {
			t.Errorf("URL %q: got %s/%s, want %s/%s", tc.url, owner, repo, tc.wantOwner, tc.wantRepo)
		}
	}
}

func TestParseRemoteURL_SSH(t *testing.T) {
	cases := []struct {
		url         string
		wantOwner   string
		wantRepo    string
	}{
		{"git@github.com:owner/repo.git", "owner", "repo"},
		{"git@github.com:owner/repo", "owner", "repo"},
		{"git@github.com:myorg/myproject.git", "myorg", "myproject"},
	}
	for _, tc := range cases {
		owner, repo, err := parseRemoteURL(tc.url)
		if err != nil {
			t.Errorf("URL %q でエラーが発生しました: %v", tc.url, err)
			continue
		}
		if owner != tc.wantOwner || repo != tc.wantRepo {
			t.Errorf("URL %q: got %s/%s, want %s/%s", tc.url, owner, repo, tc.wantOwner, tc.wantRepo)
		}
	}
}

func TestParseRemoteURL_Invalid(t *testing.T) {
	cases := []string{
		"not-a-url",
		"",
	}
	for _, url := range cases {
		_, _, err := parseRemoteURL(url)
		if err == nil {
			t.Errorf("不正な URL %q でエラーが返されるべきです", url)
		}
	}
}

// ============================================================
// filter.go のテスト
// ============================================================

func TestIsBot_BotType(t *testing.T) {
	user := User{Login: "someapp", Type: "Bot"}
	if !isBot(user) {
		t.Error("user.Type=Bot はボットと判定されるべきです")
	}
}

func TestIsBot_BotSuffix(t *testing.T) {
	cases := []string{
		"github-actions[bot]",
		"dependabot[bot]",
		"renovate[bot]",
		"app[bot]",
	}
	for _, login := range cases {
		user := User{Login: login, Type: "User"}
		if !isBot(user) {
			t.Errorf("%q は [bot] サフィックスのためボットと判定されるべきです", login)
		}
	}
}

func TestIsBot_KnownBot(t *testing.T) {
	cases := []string{"codecov", "sonarcloud", "renovate", "dependabot", "snyk-bot", "imgbot"}
	for _, login := range cases {
		user := User{Login: login, Type: "User"}
		if !isBot(user) {
			t.Errorf("%q は既知ボットとして判定されるべきです", login)
		}
	}
}

func TestIsBot_HumanUser(t *testing.T) {
	cases := []User{
		{Login: "alice", Type: "User"},
		{Login: "bob-dev", Type: "User"},
		{Login: "company-org", Type: "Organization"},
	}
	for _, user := range cases {
		if isBot(user) {
			t.Errorf("%q はボットと判定されるべきではありません", user.Login)
		}
	}
}

func TestFilterBotComments_RemovesBots(t *testing.T) {
	comments := []Comment{
		{ID: 1, Body: "人間のコメント", User: User{Login: "alice", Type: "User"}},
		{ID: 2, Body: "ボットのコメント", User: User{Login: "github-actions[bot]", Type: "Bot"}},
		{ID: 3, Body: "別の人間のコメント", User: User{Login: "bob", Type: "User"}},
		{ID: 4, Body: "codecovのコメント", User: User{Login: "codecov", Type: "User"}},
	}

	result := filterBotComments(comments)

	if len(result) != 2 {
		t.Errorf("フィルタ後のコメント数は2件のはずですが %d件 でした", len(result))
	}
	for _, c := range result {
		if isBot(c.User) {
			t.Errorf("ボットコメントが残っています: @%s", c.User.Login)
		}
	}
}

func TestFilterBotComments_EmptyInput(t *testing.T) {
	result := filterBotComments([]Comment{})
	if len(result) != 0 {
		t.Error("空のコメントリストに対して空が返されるべきです")
	}
}

func TestFilterBotComments_AllHumans(t *testing.T) {
	comments := []Comment{
		{ID: 1, Body: "コメント1", User: User{Login: "alice", Type: "User"}},
		{ID: 2, Body: "コメント2", User: User{Login: "bob", Type: "User"}},
	}
	result := filterBotComments(comments)
	if len(result) != 2 {
		t.Errorf("全て人間のコメントは全件残るべきですが %d件 でした", len(result))
	}
}

// ============================================================
// format.go のテスト
// ============================================================

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

	checks := []struct{ label, want string }{
		{"PR番号", "PR #123"},
		{"タイトル", "テスト修正"},
		{"ディスクリプション", "バグを修正しました"},
		{"マージ日", "2024-03-01"},
		{"レビュアー名", "@reviewer"},
		{"レビューコメント", "LGTMです"},
	}
	for _, c := range checks {
		if !strings.Contains(result, c.want) {
			t.Errorf("%s が含まれるべきです:\n%s", c.label, result)
		}
	}
}

func TestFormatOutput_MultiplePRs(t *testing.T) {
	details := []PRDetails{
		{PR: PullRequest{Number: 1, Title: "1つ目のPR", State: "closed"}},
		{PR: PullRequest{Number: 2, Title: "2つ目のPR", State: "open"}},
	}

	result := formatOutput(details)

	if !strings.Contains(result, "PR #1") {
		t.Errorf("1つ目のPR番号が含まれるべきです:\n%s", result)
	}
	if !strings.Contains(result, "PR #2") {
		t.Errorf("2つ目のPR番号が含まれるべきです:\n%s", result)
	}
	if !strings.Contains(result, "---") {
		t.Errorf("複数PRの区切りが含まれるべきです:\n%s", result)
	}
}

func TestFormatOutput_NoComments(t *testing.T) {
	details := []PRDetails{
		{
			PR:             PullRequest{Number: 5, Title: "コメントなしPR", Body: "変更の説明", State: "merged"},
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
		{PR: PullRequest{Number: 7, Title: "本文なしPR", Body: "", State: "closed"}},
	}

	result := formatOutput(details)

	if strings.Contains(result, "ディスクリプション") {
		t.Errorf("本文が空の場合、ディスクリプションセクションは表示されるべきでありません:\n%s", result)
	}
}

func TestFormatOutput_MergedAtTruncated(t *testing.T) {
	details := []PRDetails{
		{PR: PullRequest{Number: 10, Title: "マージ済みPR", State: "closed", MergedAt: "2024-06-15T08:30:00Z"}},
	}

	result := formatOutput(details)

	if strings.Contains(result, "T08:30:00Z") {
		t.Errorf("時刻部分は表示されるべきでありません:\n%s", result)
	}
	if !strings.Contains(result, "2024-06-15") {
		t.Errorf("日付部分が含まれるべきです:\n%s", result)
	}
}

// ============================================================
// service.go のテスト（モック使用）
// ============================================================

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
	checks := []struct{ label, want string }{
		{"PR番号", "PR #42"},
		{"タイトル", "バグ修正"},
		{"ディスクリプション", "nilポインタ参照を修正しました"},
		{"レビューコメント", "良い修正です"},
	}
	for _, c := range checks {
		if !strings.Contains(result, c.want) {
			t.Errorf("%s が含まれていません:\n%s", c.label, result)
		}
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
	// repo を明示指定した場合は RemoteRepo が呼ばれないことを確認
	git := &mockGitClient{
		sha:       "abc123def456abc123def456abc123def456abc1",
		remoteErr: errors.New("RemoteRepo が呼ばれるべきではありません"),
	}
	svc := newService(git, &mockGitHubClient{prs: []PullRequest{}})

	result, err := svc.GetPRContext("main.go", 1, "myorg/myrepo", "/tmp/repo")
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
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
