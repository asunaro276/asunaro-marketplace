package main

import (
	"fmt"
	"os"
	"strings"
)

// --- 共通型定義 ---

type User struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type PullRequest struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	State    string `json:"state"`
	MergedAt string `json:"merged_at"`
	HTMLURL  string `json:"html_url"`
}

type Comment struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	User User   `json:"user"`
}

type PRDetails struct {
	PR             PullRequest
	ReviewComments []Comment
	IssueComments  []Comment
}

// --- インターフェース定義 ---

// GitClient は git 操作を抽象化するインターフェース。
// 実装を差し替えることでテスト時にモックを使用できる。
type GitClient interface {
	BlameCommit(cwd, filePath string, line int) (string, error)
	RemoteRepo(cwd string) (owner, repo string, err error)
}

// GitHubClient は GitHub API アクセスを抽象化するインターフェース。
// 実装を差し替えることで HTTP クライアントの交換やテスト時のモックが可能。
type GitHubClient interface {
	GetPRsForCommit(owner, repo, sha string) ([]PullRequest, error)
	GetPRDetails(owner, repo string, number int) (PRDetails, error)
}

// --- サービス層 ---

type Service struct {
	git    GitClient
	github GitHubClient
}

func newService(git GitClient, github GitHubClient) *Service {
	return &Service{git: git, github: github}
}

// GetPRContext はファイルの特定行に関連する PR のコンテキストを返す。
func (s *Service) GetPRContext(filePath string, line int, repoArg, cwd string) (string, error) {
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("カレントディレクトリの取得に失敗しました: %w", err)
		}
	}

	// 1. git blame でコミット SHA を取得
	sha, err := s.git.BlameCommit(cwd, filePath, line)
	if err != nil {
		return "", fmt.Errorf("git blame の実行に失敗しました: %w", err)
	}

	// 未コミットの行（SHA が 000... で始まる）
	if strings.HasPrefix(sha, "0000000") {
		return "このコードはまだコミットされていないため、関連する PR はありません", nil
	}

	// 2. owner/repo を解決
	var owner, repo string
	if repoArg != "" {
		parts := strings.SplitN(repoArg, "/", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("repo は owner/repo の形式で指定してください")
		}
		owner, repo = parts[0], parts[1]
	} else {
		owner, repo, err = s.git.RemoteRepo(cwd)
		if err != nil {
			return "", fmt.Errorf("git remote の取得に失敗しました: %w", err)
		}
	}

	// 3. コミットに関連する PR を取得
	prs, err := s.github.GetPRsForCommit(owner, repo, sha)
	if err != nil {
		return "", fmt.Errorf("PR 一覧の取得に失敗しました: %w", err)
	}

	if len(prs) == 0 {
		shortSHA := sha
		if len(sha) > 8 {
			shortSHA = sha[:8]
		}
		return fmt.Sprintf("コミット %s は直接コミットのため、関連する PR はありません", shortSHA), nil
	}

	// 4. 各 PR の詳細を取得しボットコメントを除去
	var details []PRDetails
	for _, pr := range prs {
		d, err := s.github.GetPRDetails(owner, repo, pr.Number)
		if err != nil {
			// 1件取得できなくても他の PR の処理を継続する
			continue
		}
		d.ReviewComments = filterBotComments(d.ReviewComments)
		d.IssueComments = filterBotComments(d.IssueComments)
		details = append(details, d)
	}

	return formatOutput(details), nil
}
