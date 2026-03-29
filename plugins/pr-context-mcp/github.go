package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// httpGitHubClient は net/http を使った GitHubClient の実装。
// HTTP クライアントを差し替えたい場合はこのファイルのみ修正する。
type httpGitHubClient struct {
	token   string
	baseURL string
	client  *http.Client
}

func newGitHubClient(token string) *httpGitHubClient {
	return &httpGitHubClient{
		token:   token,
		baseURL: "https://api.github.com",
		client:  &http.Client{},
	}
}

func (c *httpGitHubClient) get(url string, v any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return fmt.Errorf("GitHub API 認証エラー (HTTP %d): GITHUB_TOKEN を確認してください", resp.StatusCode)
		}
		return fmt.Errorf("GitHub API エラー (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, v)
}

func (c *httpGitHubClient) GetPRsForCommit(owner, repo, sha string) ([]PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits/%s/pulls", c.baseURL, owner, repo, sha)
	var prs []PullRequest
	if err := c.get(url, &prs); err != nil {
		return nil, err
	}
	return prs, nil
}

func (c *httpGitHubClient) GetPRDetails(owner, repo string, number int) (PRDetails, error) {
	var details PRDetails

	// PR 本体
	prURL := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.baseURL, owner, repo, number)
	if err := c.get(prURL, &details.PR); err != nil {
		return details, err
	}

	// レビューコメント（コードへのインラインコメント）
	reviewURL := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/comments", c.baseURL, owner, repo, number)
	if err := c.get(reviewURL, &details.ReviewComments); err != nil {
		return details, err
	}

	// 一般コメント（PR スレッドのコメント）
	issueURL := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", c.baseURL, owner, repo, number)
	if err := c.get(issueURL, &details.IssueComments); err != nil {
		return details, err
	}

	return details, nil
}
