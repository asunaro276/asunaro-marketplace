package main

import "strings"

// knownBots は user.Type が "Bot" でない場合でもボットと判定する既知アカウントのセット。
var knownBots = map[string]bool{
	"codecov":    true,
	"sonarcloud": true,
	"renovate":   true,
	"dependabot": true,
	"snyk-bot":   true,
	"imgbot":     true,
	"allcontributors": true,
}

// isBot はコメント投稿者がボットかどうかを判定する。
// 判定基準:
//  1. user.Type == "Bot"（GitHub が公式にBotとマーク）
//  2. user.Login が "[bot]" で終わる（github-actions[bot] 等）
//  3. user.Login が既知ボット名リストに含まれる
func isBot(user User) bool {
	if user.Type == "Bot" {
		return true
	}
	login := strings.ToLower(user.Login)
	if strings.HasSuffix(login, "[bot]") {
		return true
	}
	return knownBots[login]
}

// filterBotComments はボットによるコメントを除外して返す。
func filterBotComments(comments []Comment) []Comment {
	result := make([]Comment, 0, len(comments))
	for _, c := range comments {
		if !isBot(c.User) {
			result = append(result, c)
		}
	}
	return result
}
