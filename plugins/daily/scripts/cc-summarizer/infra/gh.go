package infra

import (
	"cc-summarizer/domain"
	"encoding/json"
	"os/exec"
	"strings"
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
