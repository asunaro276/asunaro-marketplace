---
name: cc-daily-report:daily-report
description: このスキルは「日報を作って」「今日のClaude Code利用状況をまとめて」「デイリーレポート生成」「日報生成」「今日の作業まとめ」などのリクエストで使用します。cc-summarizerスクリプトでセッション履歴を並列解析し、Claudeが日本語の日報に仕上げます。
---

# Claude Code 日報生成スキル

## 概要

`cc-summarizer` Go スクリプトで `~/.claude/projects/` 以下のセッションファイルを並列解析し、
トークン効率よくサマリーJSONを生成する。そのJSONをClaudeが解釈して日本語の日報を作成する。

## 実行手順

### ステップ1：バイナリの準備

スクリプトのパスを特定する：

```
${CLAUDE_PLUGIN_ROOT}/scripts/cc-summarizer/main.go
```

バイナリが存在するか確認し、なければコンパイルする。OSに応じてバイナリ名を切り替える：

- Linux/macOS: `cc-summarizer`
- Windows: `cc-summarizer.exe`

バイナリが存在しない場合、以下のコマンドでコンパイルする：

```bash
cd "${CLAUDE_PLUGIN_ROOT}/scripts/cc-summarizer" && go build -o cc-summarizer ./main.go
```

Windowsの場合：
```powershell
cd "${CLAUDE_PLUGIN_ROOT}\scripts\cc-summarizer"; go build -o cc-summarizer.exe .\main.go
```

### ステップ2：スクリプト実行

引数なしで実行すると「今日（UTC日付）」のセッションを対象とする。
特定日付を指定する場合は `YYYY-MM-DD` 形式で引数を渡す：

```bash
"${CLAUDE_PLUGIN_ROOT}/scripts/cc-summarizer/cc-summarizer"
# または日付指定
"${CLAUDE_PLUGIN_ROOT}/scripts/cc-summarizer/cc-summarizer" 2026-03-28
```

出力はJSONで以下の構造を持つ（`references/report-format.md` 参照）。

### ステップ3：日報生成

JSONの内容を解釈して日本語の日報を作成する。フォーマットは `references/report-format.md` に従う。

以下の観点で内容を補完・解釈する：
- `first_user_message` からセッションのテーマ・目的を推定する
- `tools_used` からどんな作業をしたか（コーディング、検索、ファイル操作など）を推定する
- `project` フィールドからどのプロジェクトの作業かを推定する（パス名を人間が読みやすい形に変換する）
- トークン数から作業のボリューム感を把握する

### ステップ4：出力

日報をマークダウン形式で出力する。ユーザーがファイル保存を希望する場合は
`~/日報/日報_YYYY-MM-DD.md` に保存する（ディレクトリは自動作成）。

## エラーハンドリング

- `go` がインストールされていない場合：「Goのインストールが必要です（https://go.dev/dl/）」と案内する
- セッションが0件の場合：「本日はClaude Codeの利用記録がありません」と伝える
- コンパイルエラーの場合：エラーメッセージをそのまま表示して対処法を提案する

## 注意事項

- セッションファイルはUTCタイムスタンプで記録されているため、日本時間と最大9時間ずれる場合がある
- 機密情報（APIキー、パスワードなど）と思われる内容は日報に含めない
- `Warmup` メッセージのセッションは自動でスキップされる
