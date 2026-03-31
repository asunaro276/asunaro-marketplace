package infra

import (
	"bufio"
	"cc-summarizer/domain"
	"encoding/json"
	"os"
	"strings"
	"time"
)


// ReadSession reads a JSONL file, filters entries by targetDate, and returns a SessionSummary.
// Returns (nil, nil) if the file yields no valid session (e.g. no conversation entries).
func ReadSession(filePath, targetDate string) (*domain.SessionSummary, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []domain.RawEntry
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 10*1024*1024), 10*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry domain.RawEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if entry.Type == "file-history-snapshot" || entry.Timestamp == "" {
			continue
		}
		if !entry.IsForDate(targetDate) {
			continue
		}
		entries = append(entries, entry)
	}

	s := domain.NewSession(entries)
	if s == nil {
		return nil, nil
	}

	s.GitBranch = GitBranchFromCWD(s.CWD)
	if s.StartTime != "" {
		if t, err := time.Parse(time.RFC3339Nano, s.StartTime); err == nil {
			s.GitBranch = GetBranchAtTime(s.CWD, t, s.GitBranch)
		}
	}
	return s, nil
}
