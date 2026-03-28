package main

import "testing"

func TestIsBot_BotType(t *testing.T) {
	user := User{Login: "someapp", Type: "Bot"}
	if !isBot(user) {
		t.Error("user.Type=Bot はボットと判定されるべきです")
	}
}

func TestIsBot_BotSuffix(t *testing.T) {
	cases := []string{
		"github-actions[bot]",
		"dependabot[bot]",
		"renovate[bot]",
		"app[bot]",
	}
	for _, login := range cases {
		user := User{Login: login, Type: "User"}
		if !isBot(user) {
			t.Errorf("%q は [bot] サフィックスのためボットと判定されるべきです", login)
		}
	}
}

func TestIsBot_KnownBot(t *testing.T) {
	cases := []string{"codecov", "sonarcloud", "renovate", "dependabot", "snyk-bot", "imgbot"}
	for _, login := range cases {
		user := User{Login: login, Type: "User"}
		if !isBot(user) {
			t.Errorf("%q は既知ボットとして判定されるべきです", login)
		}
	}
}

func TestIsBot_HumanUser(t *testing.T) {
	cases := []User{
		{Login: "alice", Type: "User"},
		{Login: "bob-dev", Type: "User"},
		{Login: "company-org", Type: "Organization"},
	}
	for _, user := range cases {
		if isBot(user) {
			t.Errorf("%q はボットと判定されるべきではありません", user.Login)
		}
	}
}

func TestFilterBotComments_RemovesBots(t *testing.T) {
	comments := []Comment{
		{ID: 1, Body: "人間のコメント", User: User{Login: "alice", Type: "User"}},
		{ID: 2, Body: "ボットのコメント", User: User{Login: "github-actions[bot]", Type: "Bot"}},
		{ID: 3, Body: "別の人間のコメント", User: User{Login: "bob", Type: "User"}},
		{ID: 4, Body: "codecovのコメント", User: User{Login: "codecov", Type: "User"}},
	}

	result := filterBotComments(comments)

	if len(result) != 2 {
		t.Errorf("フィルタ後のコメント数は2件のはずですが %d件 でした", len(result))
	}
	for _, c := range result {
		if isBot(c.User) {
			t.Errorf("ボットコメントが残っています: @%s", c.User.Login)
		}
	}
}

func TestFilterBotComments_EmptyInput(t *testing.T) {
	result := filterBotComments([]Comment{})
	if len(result) != 0 {
		t.Errorf("空のコメントリストに対して空が返されるべきです")
	}
}

func TestFilterBotComments_AllHumans(t *testing.T) {
	comments := []Comment{
		{ID: 1, Body: "コメント1", User: User{Login: "alice", Type: "User"}},
		{ID: 2, Body: "コメント2", User: User{Login: "bob", Type: "User"}},
	}
	result := filterBotComments(comments)
	if len(result) != 2 {
		t.Errorf("全て人間のコメントは全件残るべきですが %d件 でした", len(result))
	}
}
