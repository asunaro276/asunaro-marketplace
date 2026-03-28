package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// execGitClient は os/exec を使った GitClient の実装。
// git コマンドを差し替えたい場合はこのファイルのみ修正する。
type execGitClient struct{}

func newGitClient() *execGitClient {
	return &execGitClient{}
}

func (g *execGitClient) BlameCommit(cwd, filePath string, line int) (string, error) {
	out, err := exec.Command(
		"git", "-C", cwd,
		"blame", "-L", fmt.Sprintf("%d,%d", line, line),
		"--porcelain", filePath,
	).Output()
	if err != nil {
		return "", fmt.Errorf("git blame 失敗: %w", err)
	}
	return extractSHAFromBlame(string(out))
}

// extractSHAFromBlame は git blame --porcelain の出力からコミット SHA を抽出する。
// porcelain 形式の1行目: "<40文字SHA> <orig_line> <final_line> [<num_lines>]"
func extractSHAFromBlame(output string) (string, error) {
	line := strings.SplitN(output, "\n", 2)[0]
	parts := strings.Fields(line)
	if len(parts) == 0 || len(parts[0]) < 40 {
		return "", fmt.Errorf("git blame の出力が不正です: %q", line)
	}
	return parts[0], nil
}

func (g *execGitClient) RemoteRepo(cwd string) (owner, repo string, err error) {
	out, err := exec.Command("git", "-C", cwd, "remote", "get-url", "origin").Output()
	if err != nil {
		return "", "", fmt.Errorf("git remote の取得に失敗しました: %w", err)
	}
	return parseRemoteURL(strings.TrimSpace(string(out)))
}

// parseRemoteURL は git remote の URL から owner と repo を抽出する。
// 対応フォーマット:
//   - HTTPS: https://github.com/owner/repo.git
//   - SSH:   git@github.com:owner/repo.git
func parseRemoteURL(rawURL string) (owner, repo string, err error) {
	// HTTPS / HTTP
	if strings.HasPrefix(rawURL, "https://") || strings.HasPrefix(rawURL, "http://") {
		trimmed := strings.TrimSuffix(rawURL, ".git")
		parts := strings.Split(trimmed, "/")
		if len(parts) < 2 {
			return "", "", fmt.Errorf("リモート URL の解析に失敗しました: %s", rawURL)
		}
		return parts[len(parts)-2], parts[len(parts)-1], nil
	}

	// SSH: git@github.com:owner/repo.git
	if idx := strings.Index(rawURL, ":"); idx >= 0 {
		path := strings.TrimSuffix(rawURL[idx+1:], ".git")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("リモート URL の解析に失敗しました: %s", rawURL)
		}
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("リモート URL の解析に失敗しました: %s", rawURL)
}
