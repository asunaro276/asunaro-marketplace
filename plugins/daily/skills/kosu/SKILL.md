---
name: daily:kosu
description: このスキルは「工数を出して」「PR工数確認」「今日の工数」「工数サマリー」「作業時間を出して」などのリクエストで使用します。cc-summarizerスクリプトでセッション履歴を解析し、PRごとの工数サマリーを表形式で表示します。
---

# CC 工数サマリースキル

## 概要

`cc-summarizer` Go スクリプトで `~/.claude/projects/` 以下のセッションファイルを解析し、
PRごとの作業時間（工数）をテーブル形式で表示する。

## 実行手順

### ステップ1：バイナリのパス確定

- Linux/macOS: `${CLAUDE_PLUGIN_ROOT}/scripts/cc-summarizer/cc-summarizer`
- Windows: `${CLAUDE_PLUGIN_ROOT}\scripts\cc-summarizer\cc-summarizer.exe`

ビルドが必要な場合（バイナリ実行に失敗したときのみ）：
```bash
cd "${CLAUDE_PLUGIN_ROOT}/scripts/cc-summarizer" && go build -o cc-summarizer ./main.go
```

### ステップ2：スクリプト実行

```bash
"${CLAUDE_PLUGIN_ROOT}/scripts/cc-summarizer/cc-summarizer"
# または日付指定
"${CLAUDE_PLUGIN_ROOT}/scripts/cc-summarizer/cc-summarizer" 2026-03-29
```

出力はJSON（`report/references/report-format.md` 参照）。

### ステップ3：PR工数サマリーの表示

JSONの各リポジトリの `pr_time_summary` を確認し、以下の形式で表示する。

リポジトリごとに表を出力する：

**[リポジトリ名] PR工数サマリー**

| PR | Notion チケット | 作業時間 |
|---|---|---|
| [#番号 タイトル](PR_URL) | [チケット](Notion_URL) または `-` | X.X時間 |

- `notion_url` が空の場合は `-` と表示する
- 作業時間は小数点1桁の時間単位で表示する（例: 0.8時間、1.5時間）
- `pr_time_summary` が空配列のリポジトリはスキップする

### ステップ4：合計工数の表示

全リポジトリの `pr_time_summary` の `total_hours` を合算し、最後に表示する：

```
**本日の合計工数: X.X時間**
```

全リポジトリで `pr_time_summary` が空の場合は、「今日はPRに紐づく作業記録がありません」と表示する。

## エラーハンドリング

- `go` がインストールされていない場合：「Goのインストールが必要です（https://go.dev/dl/）」と案内する
- セッションが0件の場合：「今日のセッション記録がありません」と表示する
- コンパイルエラーの場合：エラーメッセージをそのまま表示して対処法を提案する
