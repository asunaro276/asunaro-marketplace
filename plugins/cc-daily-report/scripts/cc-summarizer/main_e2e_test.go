package main

import (
	"cc-summarizer/domain"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ---- helpers ----

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "cc-summarizer")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

// makeProjectsDir creates the fake ~/.claude/projects hierarchy and returns its path.
// It also sets HOME so the binary picks it up.
func makeProjectsDir(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "projects")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// writeJSONL writes lines to <projectsDir>/<projectDirName>/<filename>.jsonl
// and fakes the file's mtime so the binary includes it for targetDate.
func writeJSONL(t *testing.T, projectsDir, projectDirName, filename, targetDate string, lines []string) {
	t.Helper()
	dir := filepath.Join(projectsDir, projectDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, filename+".jsonl")
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	// Set mtime to targetDate so the binary's ModTime filter passes.
	ts, err := time.Parse("2006-01-02", targetDate)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, ts, ts); err != nil {
		t.Fatal(err)
	}
}

// entry builders

func userEntry(sessionID, cwd, ts, text string) string {
	msg, _ := json.Marshal(text)
	return fmt.Sprintf(`{"type":"user","sessionId":%q,"timestamp":%q,"cwd":%q,"message":%s}`,
		sessionID, ts, cwd, msg)
}

func assistantEntry(sessionID, cwd, ts string, inputTokens, outputTokens int, textContent string) string {
	content := fmt.Sprintf(`[{"type":"text","text":%s}]`, mustJSON(textContent))
	msg := fmt.Sprintf(`{"usage":{"input_tokens":%d,"output_tokens":%d},"content":%s}`,
		inputTokens, outputTokens, content)
	return fmt.Sprintf(`{"type":"assistant","sessionId":%q,"timestamp":%q,"cwd":%q,"message":%s}`,
		sessionID, ts, cwd, msg)
}

func assistantWithTool(sessionID, cwd, ts string, inputTokens, outputTokens int, toolName, filePath string) string {
	input, _ := json.Marshal(map[string]string{"file_path": filePath})
	content := fmt.Sprintf(`[{"type":"tool_use","name":%s,"input":%s}]`,
		mustJSON(toolName), input)
	msg := fmt.Sprintf(`{"usage":{"input_tokens":%d,"output_tokens":%d},"content":%s}`,
		inputTokens, outputTokens, content)
	return fmt.Sprintf(`{"type":"assistant","sessionId":%q,"timestamp":%q,"cwd":%q,"message":%s}`,
		sessionID, ts, cwd, msg)
}

func mustJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func run(t *testing.T, bin string, args ...string) domain.DailyOutput {
	t.Helper()
	cmd := exec.Command(bin, args...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("binary exited with error: %v\nstdout: %s", err, out)
	}
	var result domain.DailyOutput
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	return result
}

// ---- tests ----

// UC1: 対象日のJSONLファイルが存在しない場合、空の結果を返す。
func TestNoSessions(t *testing.T) {
	bin := buildBinary(t)
	makeProjectsDir(t)

	result := run(t, bin, "2024-01-01")

	if result.Date != "2024-01-01" {
		t.Errorf("date: got %q, want 2024-01-01", result.Date)
	}
	if result.TotalRepositories != 0 {
		t.Errorf("expected 0 repositories, got %d", result.TotalRepositories)
	}
	if result.TotalSessions != 0 {
		t.Errorf("expected 0 sessions, got %d", result.TotalSessions)
	}
	if len(result.Repositories) != 0 {
		t.Errorf("expected empty repositories slice, got %v", result.Repositories)
	}
}

// UC2: 単一プロジェクト・単一セッションの基本出力構造を確認する。
func TestSingleProjectSingleSession(t *testing.T) {
	bin := buildBinary(t)
	projectsDir := makeProjectsDir(t)

	date := "2024-06-15"
	cwd := "/home/user/myapp"
	sessionID := "sess-001"
	ts1 := date + "T10:00:00Z"
	ts2 := date + "T10:05:00Z"

	writeJSONL(t, projectsDir, "-home-user-myapp", "session1", date, []string{
		userEntry(sessionID, cwd, ts1, "hello, fix the bug"),
		assistantEntry(sessionID, cwd, ts2, 100, 50, "I fixed the bug."),
	})

	result := run(t, bin, date)

	if result.Date != date {
		t.Errorf("date: got %q, want %q", result.Date, date)
	}
	if result.TotalRepositories != 1 {
		t.Errorf("total_repositories: got %d, want 1", result.TotalRepositories)
	}
	if result.TotalSessions != 1 {
		t.Errorf("total_sessions: got %d, want 1", result.TotalSessions)
	}
	if result.TotalInputTokens != 100 {
		t.Errorf("total_input_tokens: got %d, want 100", result.TotalInputTokens)
	}
	if result.TotalOutputTokens != 50 {
		t.Errorf("total_output_tokens: got %d, want 50", result.TotalOutputTokens)
	}

	repo := result.Repositories[0]
	if repo.Project != "myapp" {
		t.Errorf("project: got %q, want myapp", repo.Project)
	}
	if repo.TotalSessions != 1 {
		t.Errorf("repo total_sessions: got %d, want 1", repo.TotalSessions)
	}

	sess := repo.Sessions[0]
	if sess.SessionID != sessionID {
		t.Errorf("session_id: got %q, want %q", sess.SessionID, sessionID)
	}
	if len(sess.Conversation) != 2 {
		t.Errorf("conversation length: got %d, want 2", len(sess.Conversation))
	}
	if sess.Conversation[0].Role != "user" {
		t.Errorf("first message role: got %q, want user", sess.Conversation[0].Role)
	}
	if sess.Conversation[1].Role != "assistant" {
		t.Errorf("second message role: got %q, want assistant", sess.Conversation[1].Role)
	}
}

// UC3: 複数プロジェクトがそれぞれ独立したリポジトリとしてグルーピングされる。
func TestMultipleProjects(t *testing.T) {
	bin := buildBinary(t)
	projectsDir := makeProjectsDir(t)

	date := "2024-06-15"

	writeJSONL(t, projectsDir, "-home-user-alpha", "s1", date, []string{
		userEntry("s1", "/home/user/alpha", date+"T09:00:00Z", "fix alpha"),
		assistantEntry("s1", "/home/user/alpha", date+"T09:01:00Z", 10, 5, "done"),
	})
	writeJSONL(t, projectsDir, "-home-user-beta", "s2", date, []string{
		userEntry("s2", "/home/user/beta", date+"T10:00:00Z", "fix beta"),
		assistantEntry("s2", "/home/user/beta", date+"T10:01:00Z", 20, 8, "done"),
	})

	result := run(t, bin, date)

	if result.TotalRepositories != 2 {
		t.Errorf("total_repositories: got %d, want 2", result.TotalRepositories)
	}
	if result.TotalSessions != 2 {
		t.Errorf("total_sessions: got %d, want 2", result.TotalSessions)
	}
	if result.TotalInputTokens != 30 {
		t.Errorf("total_input_tokens: got %d, want 30", result.TotalInputTokens)
	}

	projects := make(map[string]bool)
	for _, r := range result.Repositories {
		projects[r.Project] = true
	}
	if !projects["alpha"] || !projects["beta"] {
		t.Errorf("expected projects alpha and beta, got %v", projects)
	}
}

// UC4: 同一プロジェクトの複数セッションのトークンと時間が正しく集計される。
func TestMultipleSessionsSameProject(t *testing.T) {
	bin := buildBinary(t)
	projectsDir := makeProjectsDir(t)

	date := "2024-06-15"
	cwd := "/home/user/myapp"

	writeJSONL(t, projectsDir, "-home-user-myapp", "sess-a", date, []string{
		userEntry("sess-a", cwd, date+"T09:00:00Z", "task A"),
		assistantEntry("sess-a", cwd, date+"T09:30:00Z", 100, 40, "done A"),
	})
	writeJSONL(t, projectsDir, "-home-user-myapp", "sess-b", date, []string{
		userEntry("sess-b", cwd, date+"T10:00:00Z", "task B"),
		assistantEntry("sess-b", cwd, date+"T10:20:00Z", 200, 80, "done B"),
	})

	result := run(t, bin, date)

	if result.TotalRepositories != 1 {
		t.Errorf("total_repositories: got %d, want 1 (same project should be merged)", result.TotalRepositories)
	}
	if result.TotalSessions != 2 {
		t.Errorf("total_sessions: got %d, want 2", result.TotalSessions)
	}

	repo := result.Repositories[0]
	if repo.TotalSessions != 2 {
		t.Errorf("repo total_sessions: got %d, want 2", repo.TotalSessions)
	}
	if repo.TotalInputTokens != 300 {
		t.Errorf("repo total_input_tokens: got %d, want 300", repo.TotalInputTokens)
	}
	if repo.TotalOutputTokens != 120 {
		t.Errorf("repo total_output_tokens: got %d, want 120", repo.TotalOutputTokens)
	}
	if result.TotalInputTokens != 300 {
		t.Errorf("daily total_input_tokens: got %d, want 300", result.TotalInputTokens)
	}
}

// UC5: 日付引数によるフィルタリング — 別の日のエントリは除外される。
func TestDateFiltering(t *testing.T) {
	bin := buildBinary(t)
	projectsDir := makeProjectsDir(t)

	targetDate := "2024-06-15"
	otherDate := "2024-06-14"
	cwd := "/home/user/myapp"

	// ファイルのmtimeをtargetDateにすることで取得対象に入れ、
	// 中身のtimestampだけ別日にしてフィルタリングを確認する
	writeJSONL(t, projectsDir, "-home-user-myapp", "old-session", targetDate, []string{
		userEntry("old-sess", cwd, otherDate+"T09:00:00Z", "yesterday task"),
		assistantEntry("old-sess", cwd, otherDate+"T09:01:00Z", 999, 999, "done"),
	})
	writeJSONL(t, projectsDir, "-home-user-myapp", "today-session", targetDate, []string{
		userEntry("today-sess", cwd, targetDate+"T09:00:00Z", "today task"),
		assistantEntry("today-sess", cwd, targetDate+"T09:01:00Z", 50, 20, "done today"),
	})

	result := run(t, bin, targetDate)

	if result.TotalSessions != 1 {
		t.Errorf("total_sessions: got %d, want 1 (other date should be excluded)", result.TotalSessions)
	}
	if result.TotalInputTokens != 50 {
		t.Errorf("total_input_tokens: got %d, want 50 (other date tokens must be excluded)", result.TotalInputTokens)
	}
}

// UC6: システムメッセージ (Warmup, <command-name>等) が会話履歴に含まれない。
func TestSystemMessagesExcluded(t *testing.T) {
	bin := buildBinary(t)
	projectsDir := makeProjectsDir(t)

	date := "2024-06-15"
	cwd := "/home/user/myapp"
	sessionID := "sess-sys"

	systemMessages := []string{
		"Warmup",
		"<local-command-caveat>some caveat</local-command-caveat>",
		"<command-name>commit</command-name>",
		"<local-command-stdout>output</local-command-stdout>",
		"<command-message>msg</command-message>",
		"", // empty string
	}

	lines := []string{}
	ts := date + "T10:00:00Z"
	for i, msg := range systemMessages {
		lines = append(lines, userEntry(sessionID, cwd, fmt.Sprintf("%sT10:%02d:00Z", date, i), msg))
	}
	// 通常メッセージを1件追加（セッションとして成立させるため）
	lines = append(lines, userEntry(sessionID, cwd, date+"T10:10:00Z", "normal user message"))
	lines = append(lines, assistantEntry(sessionID, cwd, ts, 10, 5, "normal reply"))

	writeJSONL(t, projectsDir, "-home-user-myapp", sessionID, date, lines)

	result := run(t, bin, date)

	if result.TotalSessions != 1 {
		t.Fatalf("expected 1 session, got %d", result.TotalSessions)
	}
	sess := result.Repositories[0].Sessions[0]
	for _, msg := range sess.Conversation {
		if msg.Role != "user" {
			continue
		}
		for _, sysMsg := range systemMessages {
			if sysMsg != "" && msg.Content == sysMsg {
				t.Errorf("system message leaked into conversation: %q", msg.Content)
			}
		}
		if msg.Content == "" {
			t.Errorf("empty message leaked into conversation")
		}
	}
}

// UC7: ツール使用が tools_used に集計され、アクセスファイルが files_accessed に記録される。
func TestToolsAndFilesTracking(t *testing.T) {
	bin := buildBinary(t)
	projectsDir := makeProjectsDir(t)

	date := "2024-06-15"
	cwd := "/home/user/myapp"
	sessionID := "sess-tools"

	writeJSONL(t, projectsDir, "-home-user-myapp", sessionID, date, []string{
		userEntry(sessionID, cwd, date+"T10:00:00Z", "refactor the code"),
		assistantWithTool(sessionID, cwd, date+"T10:01:00Z", 50, 20, "Read", "/home/user/myapp/src/main.go"),
		assistantWithTool(sessionID, cwd, date+"T10:02:00Z", 50, 20, "Edit", "/home/user/myapp/src/main.go"),
		assistantWithTool(sessionID, cwd, date+"T10:03:00Z", 50, 20, "Read", "/home/user/myapp/src/util.go"),
	})

	result := run(t, bin, date)

	if result.TotalSessions != 1 {
		t.Fatalf("expected 1 session, got %d", result.TotalSessions)
	}
	sess := result.Repositories[0].Sessions[0]

	if sess.ToolsUsed["Read"] != 2 {
		t.Errorf("Read tool count: got %d, want 2", sess.ToolsUsed["Read"])
	}
	if sess.ToolsUsed["Edit"] != 1 {
		t.Errorf("Edit tool count: got %d, want 1", sess.ToolsUsed["Edit"])
	}

	fileSet := make(map[string]bool)
	for _, f := range sess.FilesAccessed {
		fileSet[f] = true
	}
	if !fileSet["/home/user/myapp/src/main.go"] {
		t.Error("main.go not found in files_accessed")
	}
	if !fileSet["/home/user/myapp/src/util.go"] {
		t.Error("util.go not found in files_accessed")
	}
}

// UC8: セッション時間 (duration_minutes) が start/end timestamp から正しく計算される。
func TestSessionDuration(t *testing.T) {
	bin := buildBinary(t)
	projectsDir := makeProjectsDir(t)

	date := "2024-06-15"
	cwd := "/home/user/myapp"
	sessionID := "sess-dur"

	writeJSONL(t, projectsDir, "-home-user-myapp", sessionID, date, []string{
		userEntry(sessionID, cwd, date+"T10:00:00Z", "start work"),
		assistantEntry(sessionID, cwd, date+"T10:30:00Z", 10, 5, "done after 30 min"),
	})

	result := run(t, bin, date)

	if result.TotalSessions != 1 {
		t.Fatalf("expected 1 session, got %d", result.TotalSessions)
	}
	sess := result.Repositories[0].Sessions[0]
	if sess.DurationMinutes != 30.0 {
		t.Errorf("duration_minutes: got %v, want 30.0", sess.DurationMinutes)
	}
}

// UC9: file-history-snapshot エントリは無視される。
func TestFileHistorySnapshotIgnored(t *testing.T) {
	bin := buildBinary(t)
	projectsDir := makeProjectsDir(t)

	date := "2024-06-15"
	cwd := "/home/user/myapp"
	sessionID := "sess-snap"

	snapshot := fmt.Sprintf(`{"type":"file-history-snapshot","sessionId":%q,"timestamp":%q,"cwd":%q,"message":{}}`,
		sessionID, date+"T09:00:00Z", cwd)

	writeJSONL(t, projectsDir, "-home-user-myapp", sessionID, date, []string{
		snapshot,
		userEntry(sessionID, cwd, date+"T10:00:00Z", "actual work"),
		assistantEntry(sessionID, cwd, date+"T10:01:00Z", 10, 5, "done"),
	})

	result := run(t, bin, date)

	// snapshot が無視されても session は成立する
	if result.TotalSessions != 1 {
		t.Errorf("expected 1 session, got %d", result.TotalSessions)
	}
}
