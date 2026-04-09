package infra

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GitBranchFromCWD(cwd string) string {
	if cwd == "" {
		return ""
	}
	// cwdから親ディレクトリを辿って .git を探す（サブディレクトリ・worktree対応）
	dir := cwd
	for {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			if info.IsDir() {
				// 通常のgitリポジトリ
				return readHeadFile(filepath.Join(gitPath, "HEAD"))
			}
			// git worktree: .git はファイルで "gitdir: /path/to/.git/worktrees/xxx" を含む
			data, err := os.ReadFile(gitPath)
			if err != nil {
				return ""
			}
			line := strings.TrimSpace(string(data))
			if !strings.HasPrefix(line, "gitdir: ") {
				return ""
			}
			gitdir := strings.TrimPrefix(line, "gitdir: ")
			if !filepath.IsAbs(gitdir) {
				gitdir = filepath.Join(dir, gitdir)
			}
			return readHeadFile(filepath.Join(gitdir, "HEAD"))
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func readHeadFile(headPath string) string {
	data, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	if strings.HasPrefix(line, "ref: refs/heads/") {
		return strings.TrimPrefix(line, "ref: refs/heads/")
	}
	if len(line) >= 7 {
		return line[:7]
	}
	return line
}

// checkoutEvent はreflogのcheckoutイベントを表す。
type checkoutEvent struct {
	at   time.Time
	from string
	to   string
}

// parseCheckoutEvents はgit reflog出力からcheckoutイベントを新しい順にパースする。
func parseCheckoutEvents(reflogOutput string) []checkoutEvent {
	var events []checkoutEvent
	for _, line := range strings.Split(strings.TrimSpace(reflogOutput), "\n") {
		if !strings.Contains(line, "checkout: moving from ") {
			continue
		}
		// %ci format: "2006-01-02 15:04:05 -0700"
		const dtLen = len("2006-01-02 15:04:05 -0700")
		if len(line) <= dtLen {
			continue
		}
		eventTime, err := time.Parse("2006-01-02 15:04:05 -0700", line[:dtLen])
		if err != nil {
			continue
		}
		rest := line[dtLen+1:] // "checkout: moving from feat/xxx to main"
		idx := strings.LastIndex(rest, " to ")
		if idx < 0 {
			continue
		}
		// "checkout: moving from " の後から " to " の前までがfromブランチ
		fromPart := rest[len("checkout: moving from "):idx]
		fromBranch := strings.TrimSpace(fromPart)
		toBranch := strings.TrimSpace(rest[idx+4:])
		events = append(events, checkoutEvent{eventTime, fromBranch, toBranch})
	}
	return events
}

// resolveBranchAtTime はcheckoutイベント列（新しい順）から時刻tでのブランチを判定する。
func resolveBranchAtTime(events []checkoutEvent, t time.Time, fallback string) string {
	// events are newest-first; find the latest event at or before t
	for _, e := range events {
		if !e.at.After(t) {
			return e.to
		}
	}

	// セッション開始時刻より前のcheckoutイベントがない場合、
	// 最も古いcheckoutイベントの「from」ブランチを返す。
	// 最古のcheckoutが「from A to B」なら、それ以前はブランチAだったはず。
	if len(events) > 0 {
		return events[len(events)-1].from
	}
	return fallback
}

// GetBranchAtTime returns the git branch that was checked out at time t
// by parsing git reflog checkout entries. Returns fallback if not determinable.
func GetBranchAtTime(cwd string, t time.Time, fallback string) string {
	cmd := exec.Command("git", "reflog", "--format=%ci %gs")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return fallback
	}

	events := parseCheckoutEvents(string(out))
	return resolveBranchAtTime(events, t, fallback)
}
