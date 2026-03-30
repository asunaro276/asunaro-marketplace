package infra

import (
	"cc-summarizer/domain"
	"encoding/json"
	"os/exec"
)

func FetchPRInfo(cwd, branch string) *domain.PRInfo {
	if branch == "" || branch == "main" || branch == "master" {
		return nil
	}
	cmd := exec.Command("gh", "pr", "list",
		"--head", branch,
		"--state", "all",
		"--json", "number,title,url,body",
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
	}
	if json.Unmarshal(out, &prs) != nil || len(prs) == 0 {
		return nil
	}
	pr := prs[0]
	return domain.NewPRInfo(pr.Number, pr.Title, pr.URL, pr.Body)
}
