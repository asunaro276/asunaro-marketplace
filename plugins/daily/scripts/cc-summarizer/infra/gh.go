package infra

import (
	"cc-summarizer/domain"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// FetchGitHubUsername は現在認証中の GitHub ユーザー名を返す。
// 取得に失敗した場合は空文字を返す。
func FetchGitHubUsername() string {
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func FetchPRInfo(cwd, branch string) *domain.PRInfo {
	if branch == "" || branch == "main" || branch == "master" {
		return nil
	}
	cmd := exec.Command("gh", "pr", "list",
		"--head", branch,
		"--state", "all",
		"--json", "number,title,url,body,author",
		"--limit", "1",
	)
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	var prs []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		URL    string `json:"url"`
		Body   string `json:"body"`
		Author struct {
			Login string `json:"login"`
		} `json:"author"`
	}
	if json.Unmarshal(out, &prs) != nil || len(prs) == 0 {
		return nil
	}
	pr := prs[0]
	return domain.NewPRInfo(pr.Number, pr.Title, pr.URL, pr.Body, pr.Author.Login)
}

// FetchPRInfoByCommits はセッション時間帯のコミットからPRを逆引きする。
// ブランチベースのPR解決が失敗した場合のフォールバックとして使用する。
// git name-rev を使ってコミットの最も近いブランチを高速に特定する。
func FetchPRInfoByCommits(cwd string, startTime, endTime time.Time, ghUser string) *domain.PRInfo {
	if cwd == "" {
		return nil
	}

	// セッション中のコミットSHAを取得（著者フィルタ付き）
	args := []string{
		"log",
		"--after=" + startTime.Add(-1*time.Minute).Format(time.RFC3339),
		"--before=" + endTime.Add(1*time.Minute).Format(time.RFC3339),
		"--format=%H",
		"--all",
	}
	if ghUser != "" {
		args = append(args, fmt.Sprintf("--author=%s", ghUser))
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	shas := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(shas) == 0 || shas[0] == "" {
		return nil
	}

	// 最大5件のコミットだけチェック
	if len(shas) > 5 {
		shas = shas[:5]
	}

	// git name-rev でコミットの最も近いブランチを特定し、PRを検索
	seen := make(map[string]bool)
	for _, sha := range shas {
		nrCmd := exec.Command("git", "name-rev",
			"--refs=refs/remotes/origin/*",
			"--exclude=refs/remotes/origin/main",
			"--exclude=refs/remotes/origin/master",
			"--exclude=refs/remotes/origin/staging",
			"--exclude=refs/remotes/origin/develop",
			"--exclude=refs/remotes/origin/HEAD",
			sha,
		)
		nrCmd.Dir = cwd
		nrOut, err := nrCmd.Output()
		if err != nil {
			continue
		}

		// 出力形式: "<sha> remotes/origin/<branch>~N" or "<sha> remotes/origin/<branch>"
		line := strings.TrimSpace(string(nrOut))
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 || parts[1] == "undefined" {
			continue
		}
		ref := parts[1]
		// "~N" や "^N" サフィックスを除去
		for _, sep := range []string{"~", "^"} {
			if idx := strings.Index(ref, sep); idx >= 0 {
				ref = ref[:idx]
			}
		}
		// "remotes/origin/" プレフィックスを除去
		branch := strings.TrimPrefix(ref, "remotes/origin/")
		if branch == ref || branch == "" {
			continue
		}
		if seen[branch] {
			continue
		}
		seen[branch] = true

		prInfo := FetchPRInfo(cwd, branch)
		if prInfo != nil {
			return prInfo
		}
	}
	return nil
}
