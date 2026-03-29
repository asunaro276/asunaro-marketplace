# 日報フォーマット仕様

## cc-summarizer の出力JSON構造

```json
{
  "date": "2026-03-29",
  "total_repositories": 2,
  "total_sessions": 10,
  "total_input_tokens": 45000,
  "total_output_tokens": 120000,
  "repositories": [
    {
      "project": "myapp",
      "total_sessions": 4,
      "total_input_tokens": 20000,
      "total_output_tokens": 60000,
      "start_time": "2026-03-29T08:12:00.000Z",
      "end_time": "2026-03-29T10:30:00.000Z",
      "git_branches": ["main", "feature/login"],
      "message_count": 48,
      "turn_count": 24,
      "tools_used": {"Bash": 12, "Read": 8, "Edit": 5, "Agent": 2},
      "files_accessed": [
        "/home/user/myapp/src/auth.ts",
        "/home/user/myapp/src/login.vue"
      ],
      "sessions": [
        {
          "session_id": "uuid-...",
          "project": "myapp",
          "git_branch": "main",
          "start_time": "2026-03-29T08:12:00.000Z",
          "end_time": "2026-03-29T09:45:00.000Z",
          "user_prompts": [
            "ログイン機能のバグを修正してほしい",
            "テストも追加して"
          ],
          "assistant_responses": [
            "ログイン処理のバグを特定しました。auth.ts の42行目で...",
            "テストを追加しました。以下の3ケースをカバーしています..."
          ],
          "total_input_tokens": 15000,
          "total_output_tokens": 40000,
          "message_count": 24,
          "turn_count": 12,
          "tools_used": {"Bash": 5, "Read": 4, "Edit": 3},
          "files_accessed": [
            "/home/user/myapp/src/auth.ts"
          ]
        }
      ]
    }
  ]
}
```

### フィールド解説

#### DailyOutput（トップレベル）

| フィールド | 説明 |
|---|---|
| `date` | 対象日（UTC） |
| `total_repositories` | 当日に作業したGitリポジトリ数 |
| `total_sessions` | 全リポジトリのセッション数合計 |
| `total_input_tokens` | 全セッションの入力トークン合計 |
| `total_output_tokens` | 全セッションの出力トークン合計 |

#### RepositorySummary

| フィールド | 説明 |
|---|---|
| `project` | Gitルートディレクトリ名（`git rev-parse --show-toplevel` の basename） |
| `total_sessions` | このリポジトリのセッション数 |
| `git_branches` | 使用したブランチ一覧（重複排除・ソート済み） |
| `message_count` | user + assistant メッセージ数合計 |
| `turn_count` | ユーザーのターン数合計 |
| `tools_used` | ツール名→使用回数のマップ（全セッション集計） |
| `files_accessed` | 操作したファイル・ディレクトリ一覧（重複排除・ソート済み） |

#### SessionSummary

| フィールド | 説明 |
|---|---|
| `session_id` | セッションの一意識別子 |
| `project` | リポジトリ名 |
| `git_branch` | セッション開始時のGitブランチ |
| `start_time` / `end_time` | セッションの開始・終了時刻（UTC） |
| `user_prompts` | ユーザーが入力した全プロンプト（各200文字上限） |
| `assistant_responses` | Claudeの出力テキスト（各300文字上限） |
| `total_input_tokens` | このセッションの入力トークン数 |
| `total_output_tokens` | このセッションの出力トークン数 |
| `message_count` | user + assistant メッセージ数 |
| `turn_count` | ユーザーのターン数 |
| `tools_used` | ツール名→使用回数のマップ |
| `files_accessed` | 操作したファイル一覧（Read/Write/Edit/Glob/Grep から抽出） |

### プロジェクト名の解読

`project` フィールドは `git rev-parse --show-toplevel` の出力の basename から取得します。例：
- `/home/ryuhei/asunaro-marketplace` → `asunaro-marketplace`
- Gitリポジトリでない場合は作業ディレクトリ（`cwd`）の basename

git worktree で作業した場合も、同一リポジトリのルートにまとめられます。

### ツール名から作業内容を推定

| ツール | 示す作業 |
|---|---|
| `Read`, `Glob`, `Grep` | コード調査・検索 |
| `Edit`, `Write` | コーディング・ファイル編集 |
| `Bash` | コマンド実行・テスト・ビルド |
| `WebSearch`, `WebFetch` | 調査・情報収集 |
| `Agent` | 複雑な自律タスク |
| `mcp__*` | 外部サービス連携 |

## Daily/への書き込みフォーマット

`~/zettelkasten/Daily/YYYY-MM-DD.md` の `# ふりかえり` セクションに書き込む内容：

```markdown
# ふりかえり

## Claude Code作業ログ
- **[プロジェクト名]**: [ユーザープロンプトとツール使用から読み取れる作業内容を事実として列挙]
- **[プロジェクト名]**: ...

## [振り返り軸タイトル]

**Q: [問い]**
[ユーザーの回答。要約・言い換えをしない。フォローアップがあった場合は改行して続ける]

**Q: [問い]**
[ユーザーの回答]

...
```

### 書き込み仕様

- `# ふりかえり` 以降をすべて置き換える
- `# メモ` セクションは変更しない
- ユーザーの回答は語尾・言い回し含めてそのまま記録する（要約・整形禁止）
- Claude Code作業ログは「〜しました」などの評価を加えず事実のみ記述する

## 注意事項

- タイムスタンプはUTC。日本時間（JST）に変換するには+9時間する
- `user_prompts` が空の場合はそのセッションをスキップ
- リポジトリが多い場合（5件以上）は類似テーマをグループ化してよい
- 機密情報（APIキー、パスワードなど）と思われる内容は含めない
