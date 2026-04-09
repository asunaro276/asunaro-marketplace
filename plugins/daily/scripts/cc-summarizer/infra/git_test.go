package infra

import (
	"testing"
	"time"
)

func mustParseTime(t *testing.T, s string) time.Time {
	t.Helper()
	tm, err := time.Parse("2006-01-02 15:04:05 -0700", s)
	if err != nil {
		t.Fatalf("mustParseTime: %v", err)
	}
	return tm
}

func TestParseCheckoutEvents(t *testing.T) {
	reflog := `2026-04-05 14:00:00 +0900 checkout: moving from feature/b to main
2026-04-03 11:30:46 +0900 checkout: moving from feature/a to feature/b
2026-04-01 09:00:00 +0900 commit: some commit message`

	events := parseCheckoutEvents(reflog)
	if len(events) != 2 {
		t.Fatalf("イベント数が想定と異なる: got %d, want 2", len(events))
	}

	// 1つ目（最新）
	if events[0].from != "feature/b" || events[0].to != "main" {
		t.Errorf("events[0]: got from=%q to=%q, want from=feature/b to=main", events[0].from, events[0].to)
	}

	// 2つ目（最古）
	if events[1].from != "feature/a" || events[1].to != "feature/b" {
		t.Errorf("events[1]: got from=%q to=%q, want from=feature/a to=feature/b", events[1].from, events[1].to)
	}
}

func TestResolveBranchAtTime_NormalCase(t *testing.T) {
	// イベント: 4/3に feature/a → feature/b, 4/5に feature/b → main
	events := []checkoutEvent{
		{mustParseTime(t, "2026-04-05 14:00:00 +0900"), "feature/b", "main"},
		{mustParseTime(t, "2026-04-03 11:30:46 +0900"), "feature/a", "feature/b"},
	}

	tests := []struct {
		name     string
		t        time.Time
		fallback string
		want     string
	}{
		{
			name:     "4/6のセッション → main（最新イベント以降）",
			t:        mustParseTime(t, "2026-04-06 10:00:00 +0900"),
			fallback: "current",
			want:     "main",
		},
		{
			name:     "4/4のセッション → feature/b（4/3のcheckout以降、4/5以前）",
			t:        mustParseTime(t, "2026-04-04 10:00:00 +0900"),
			fallback: "current",
			want:     "feature/b",
		},
		{
			name:     "4/3 11:30:46のセッション → feature/b（ちょうどイベント時刻）",
			t:        mustParseTime(t, "2026-04-03 11:30:46 +0900"),
			fallback: "current",
			want:     "feature/b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveBranchAtTime(events, tc.t, tc.fallback)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestResolveBranchAtTime_FallbackToOldestFrom(t *testing.T) {
	// バグ修正の本題: reflogの最古イベントよりも前のセッション
	// 例: worktreeで4/2に feature/improve-lcp-performance で作業
	//     4/3にcheckout: moving from feature/improve-lcp-performance to main
	//     → 4/2のセッションは feature/improve-lcp-performance であるべき
	events := []checkoutEvent{
		{mustParseTime(t, "2026-04-05 10:00:00 +0900"), "main", "feat/sentry-triage"},
		{mustParseTime(t, "2026-04-03 11:30:46 +0900"), "feature/improve-lcp-performance", "main"},
	}

	// 4/2のセッション → 最古イベントのfromを返すべき
	got := resolveBranchAtTime(events, mustParseTime(t, "2026-04-02 15:00:00 +0900"), "current-branch-fallback")
	want := "feature/improve-lcp-performance"
	if got != want {
		t.Errorf("最古イベント以前のセッション: got %q, want %q", got, want)
	}
}

func TestResolveBranchAtTime_NoEvents(t *testing.T) {
	// checkoutイベントがない場合はfallbackを返す
	got := resolveBranchAtTime(nil, mustParseTime(t, "2026-04-02 15:00:00 +0900"), "my-fallback")
	if got != "my-fallback" {
		t.Errorf("イベントなし: got %q, want %q", got, "my-fallback")
	}
}

func TestResolveBranchAtTime_SingleEvent(t *testing.T) {
	events := []checkoutEvent{
		{mustParseTime(t, "2026-04-03 11:30:46 +0900"), "feature/a", "main"},
	}

	// イベント以前 → feature/a（fromを返す）
	got := resolveBranchAtTime(events, mustParseTime(t, "2026-04-02 10:00:00 +0900"), "fallback")
	if got != "feature/a" {
		t.Errorf("単一イベント以前: got %q, want %q", got, "feature/a")
	}

	// イベント以降 → main（toを返す）
	got = resolveBranchAtTime(events, mustParseTime(t, "2026-04-04 10:00:00 +0900"), "fallback")
	if got != "main" {
		t.Errorf("単一イベント以降: got %q, want %q", got, "main")
	}
}
