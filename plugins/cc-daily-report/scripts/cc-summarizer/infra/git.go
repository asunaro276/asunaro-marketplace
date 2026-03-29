package infra

import (
	"os"
	"path/filepath"
	"strings"
)

func GitBranchFromCWD(cwd string) string {
	if cwd == "" {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(cwd, ".git", "HEAD"))
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
