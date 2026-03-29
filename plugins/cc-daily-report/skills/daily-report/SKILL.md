---
name: cc-daily-report:daily-report
description: このスキルは「日報を作って」「今日のClaude Code利用状況をまとめて」「デイリーレポート生成」「日報生成」「今日の作業まとめ」などのリクエストで使用します。cc-summarizerスクリプトでセッション履歴を並列解析し、インタラクティブな振り返りQ&Aを経てZettelkastenのDailyノートに書き込みます。
---

# Claude Code 日報・振り返りスキル

## 概要

`cc-summarizer` Go スクリプトで `~/.claude/projects/` 以下のセッションファイルを並列解析し、
Claude Code作業データと今日のメモを素材にインタラクティブな振り返りを行う。
回答内容を `~/zettelkasten/Daily/YYYY-MM-DD.md` の `# ふりかえり` セクションに書き込む。

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

出力はJSON（`references/report-format.md` 参照）。

### ステップ3：デイリーノートの読み込み

`~/zettelkasten/Daily/YYYY-MM-DD.md` を読む。存在しない場合は空として扱う。
`# メモ` セクションの内容を振り返りの素材として使う。

### ステップ4：インタラクティブ振り返りQ&A

#### 4-1. 作業事実の提示

CC作業データから今日のプロジェクトと作業内容を **事実として** 箇条書きで提示する。
今日のメモがあればそれも合わせて表示する。
「〜を達成しました」などの評価・要約を加えない。

#### 4-2. 振り返り問いを1問ずつ聞く

**問いの構造**: 必ず **「前提 → 問い」** の2段構成にする。

```
[前提] 作業内容をuser_prompts・assistant_responses・files_accessedから1〜3文で再現する
[問い] 作業の性質に応じた問い（固有名詞を必ず含める）
```

例：
> cc-summarizerで、git worktreeも含めて同一リポジトリとみなして集計するロジックを実装しました。
> このリポジトリ集計ロジックの本質は何だと思いますか？

**問いの選び方**: 作業の性質に応じて `references/question-types.md` のテーブルを参照する。
目的が不明確なとき → ①目的を問う、目的が明確なとき → ②妥当性を問う。

**NG**:
- 答えの方向を誘導しない
- 抽象的な「今日の作業」は使わない（固有名詞を使う）

**フォローアップ**: 回答が表面的・曖昧なら1回だけ掘り下げる。

#### 4-3. 全問終了後

全問の回答が揃ったら、確認なしにそのままステップ5に進む。

### ステップ5：Daily/への書き込み

`~/zettelkasten/Daily/YYYY-MM-DD.md` の `# ふりかえり` セクションを置き換える。
書き込む内容のフォーマットは `references/report-format.md` 参照。

**仕様**:
- `# ふりかえり` 以降をすべて新しい内容に置き換える
- `# メモ` セクションは一切変更しない
- ファイルが存在しない場合は新規作成する

書き込み完了後、ファイルパスをユーザーに伝える。

### ステップ5-2：PR工数レポートの表示

JSONの各リポジトリに `pr_time_summary` が含まれている場合、以下の表をユーザーに表示する：

**[リポジトリ名] PR工数サマリー**

| PR | Notion チケット | 作業時間 |
|---|---|---|
| [#番号 タイトル](PR_URL) | [チケット](Notion_URL) または `-` | X.X時間 |

- `notion_url` が空の場合は `-` と表示する
- 作業時間は小数点1桁の時間単位で表示する（例: 0.8時間、1.5時間）
- `pr_time_summary` が空配列の場合はこのステップをスキップする

## エラーハンドリング

- `go` がインストールされていない場合：「Goのインストールが必要です（https://go.dev/dl/）」と案内する
- セッションが0件の場合：振り返りは続行する（Claude Code未使用の日もある）
- コンパイルエラーの場合：エラーメッセージをそのまま表示して対処法を提案する

## 注意事項

- セッションファイルはUTCタイムスタンプで記録されているため、日本時間と最大9時間ずれる場合がある
- 機密情報（APIキー、パスワードなど）と思われる内容は日報に含めない
- `Warmup` メッセージのセッションは自動でスキップされる
- ユーザーの回答は要約せず、語尾や言い回しも含めてそのまま記録する
