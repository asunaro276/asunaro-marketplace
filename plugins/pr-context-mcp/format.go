package main

import (
	"fmt"
	"strings"
)

// formatOutput は PR 詳細一覧を Markdown 形式の文字列にフォーマットする。
func formatOutput(details []PRDetails) string {
	if len(details) == 0 {
		return "関連する PR が見つかりませんでした"
	}

	var sb strings.Builder
	for i, d := range details {
		if i > 0 {
			sb.WriteString("\n---\n\n")
		}
		formatPR(&sb, d)
	}
	return sb.String()
}

func formatPR(sb *strings.Builder, d PRDetails) {
	pr := d.PR

	fmt.Fprintf(sb, "## PR #%d: %s\n", pr.Number, pr.Title)
	fmt.Fprintf(sb, "**状態**: %s", pr.State)
	if pr.MergedAt != "" {
		mergedDate := pr.MergedAt
		if len(mergedDate) >= 10 {
			mergedDate = mergedDate[:10]
		}
		fmt.Fprintf(sb, " (マージ日: %s)", mergedDate)
	}
	sb.WriteString("\n\n")

	if strings.TrimSpace(pr.Body) != "" {
		sb.WriteString("### ディスクリプション\n")
		sb.WriteString(pr.Body)
		sb.WriteString("\n\n")
	}

	if len(d.ReviewComments) > 0 {
		fmt.Fprintf(sb, "### レビューコメント (%d件)\n", len(d.ReviewComments))
		for _, c := range d.ReviewComments {
			fmt.Fprintf(sb, "- @%s: %s\n", c.User.Login, c.Body)
		}
		sb.WriteString("\n")
	}

	if len(d.IssueComments) > 0 {
		fmt.Fprintf(sb, "### 一般コメント (%d件)\n", len(d.IssueComments))
		for _, c := range d.IssueComments {
			fmt.Fprintf(sb, "- @%s: %s\n", c.User.Login, c.Body)
		}
		sb.WriteString("\n")
	}
}
