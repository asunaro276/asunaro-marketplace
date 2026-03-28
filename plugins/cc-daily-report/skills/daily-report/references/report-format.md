# 日報フォーマット仕様

## cc-summarizer の出力JSON構造

```json
{
  "date": "2026-03-28",
  "total_sessions": 3,
  "total_input_tokens": 45000,
  "total_output_tokens": 12000,
  "sessions": [
    {
      "session_id": "uuid-...",
      "project": "C--Users-ryuhe-projects-myapp",
      "start_time": "2026-03-28T08:12:00.000Z",
      "end_time": "2026-03-28T09:45:00.000Z",
      "first_user_message": "ログイン機能のバグを修正してほしい",
      "total_input_tokens": 15000,
      "total_output_tokens": 4000,
      "tools_used": ["Bash", "Edit", "Read"],
      "turn_count": 12
    }
  ]
}
```

### フィールド解説

| フィールド | 説明 |
|---|---|
| `date` | 対象日（UTC） |
| `total_sessions` | 当日のセッション数 |
| `total_input_tokens` | 全セッションの入力トークン合計 |
| `total_output_tokens` | 全セッションの出力トークン合計 |
| `project` | `~/.claude/projects/` 以下のフォルダ名（パスのスラッシュを`--`に変換したもの） |
| `first_user_message` | そのセッションの最初のユーザーメッセージ（最大200文字） |
| `tools_used` | 使用したツール名リスト |
| `turn_count` | ユーザーのターン数（会話の往復回数） |

### プロジェクト名の解読

`project` フィールドはパスをエンコードしたもの。例：
- `C--Users-ryuhe` → ホームディレクトリ（ルート作業）
- `C--Users-ryuhe-projects-myapp` → `~/projects/myapp`
- `C--Users-ryuhe-OneDrive--------zettelkasten` → OneDrive内のzettelkastenフォルダ

表示時は末尾のフォルダ名を使って「myappプロジェクト」のように読みやすくする。

### ツール名から作業内容を推定

| ツール | 示す作業 |
|---|---|
| `Read`, `Glob`, `Grep` | コード調査・検索 |
| `Edit`, `Write` | コーディング・ファイル編集 |
| `Bash` | コマンド実行・テスト・ビルド |
| `WebSearch`, `WebFetch` | 調査・情報収集 |
| `Agent` | 複雑な自律タスク |
| `mcp__*` | 外部サービス連携 |

## 日報テンプレート

```markdown
# 日報 YYYY年MM月DD日（曜日）

## サマリー
（1〜2文で今日全体の作業を要約）

## セッション別作業内容

### 1. [セッションのテーマ]（プロジェクト名）
- **時間帯**: HH:MM〜HH:MM（JST換算）
- **作業内容**: （first_user_message と tools_used から推定）
- **トークン使用**: 入力 X,XXX / 出力 X,XXX

### 2. ...

## 本日の統計
- セッション数: N 件
- 総トークン使用: 入力 XX,XXX / 出力 XX,XXX
- 主な作業カテゴリ: （コーディング / 調査 / 設計 など）

## 所感・メモ
（作業の振り返りや気づきがあれば）
```

## 注意事項

- タイムスタンプはUTC。日本時間（JST）に変換するには+9時間する
- `first_user_message` が空の場合は「内容不明」と記載
- セッションが多い場合（5件以上）は類似テーマをグループ化してよい
