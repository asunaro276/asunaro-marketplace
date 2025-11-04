# CLAUDE.md

このファイルは、このリポジトリのコードを扱う際に Claude Code (claude.ai/code) へのガイダンスを提供します。

## 言語
ドキュメント・標準出力は日本語で出力してください。

## リポジトリ概要

これは Claude Code 用のプライベートプラグインマーケットプレイスです。プラグインとスキルを管理し、Claude Code の機能を拡張します。

## アーキテクチャ

### ディレクトリ構造

```
asunaro-marketplace/
├── .claude-plugin/
│   └── marketplace.json        # マーケットプレイスのメタデータ
└── plugins/
    └── development-plugin/     # 開発用プラグイン
        ├── .claude-plugin/
        │   └── plugin.json     # プラグイン設定（MCP サーバー含む）
        └── skills/             # プラグインスキル
            └── git-operations/ # Git 操作スキル
                ├── SKILL.md
                └── references/
```

### マーケットプレイス構造

- **ルートレベル**: `.claude-plugin/marketplace.json` がマーケットプレイス全体を定義
- **プラグインレベル**: 各プラグインは `plugins/` ディレクトリ内に独自の `.claude-plugin/plugin.json` を持つ
- **スキルレベル**: プラグインは `skills/` ディレクトリ内に複数のスキルを含むことができる

### プラグイン設定

プラグインは以下を定義できます：
- **MCPサーバー**: `mcpServers` フィールドで外部 MCP サーバーを統合
  - 例：`development-plugin` は Serena MCP サーバー（IDEアシスタント）を統合

### スキルシステム

スキルは次の形式で構造化されます：
- `SKILL.md`: スキルのメインドキュメント（フロントマター付き）
  - `name`: スキル識別子
  - `description`: スキルをいつ使用するか
- `references/`: 詳細なリファレンスドキュメント

## 開発ワークフロー

### 新しいプラグインの追加

1. `plugins/` ディレクトリ内に新しいディレクトリを作成
2. `.claude-plugin/plugin.json` をプラグイン設定で作成
3. `.claude-plugin/marketplace.json` の `plugins` 配列に登録
4. 必要に応じて `skills/` ディレクトリを追加

### 新しいスキルの追加

1. プラグインの `skills/` ディレクトリ内にスキルディレクトリを作成
2. フロントマター（`name`、`description`）を含む `SKILL.md` を作成
3. 必要に応じて `references/` ディレクトリに詳細ドキュメントを追加

### メタデータの更新

プラグインまたはマーケットプレイスのメタデータを変更する場合：
- **プラグイン**: `plugins/<plugin-name>/.claude-plugin/plugin.json` を編集
- **マーケットプレイス**: `.claude-plugin/marketplace.json` を編集

### MCP サーバーの統合

プラグインに新しい MCP サーバーを追加するには：
1. プラグインの `plugin.json` を編集
2. `mcpServers` オブジェクトに新しいサーバー設定を追加：
   ```json
   "mcpServers": {
     "server-name": {
       "command": "command-to-run",
       "args": ["arg1", "arg2"]
     }
   }
   ```

## 既存のプラグイン

### development-plugin

開発を行う際に言語にかかわらず必要となる Claude Plugin。

**統合された MCP サーバー:**
- **Serena**: IDE アシスタントコンテキストを提供（`uvx` 経由）

**スキル:**
- **git-operations**: Git 操作に関する包括的なガイド
  - ブランチ作成、コミット、プッシュ、プル、マージ、リベース
  - Git のベストプラクティスと安全性プロトコル
  - ネットワークエラーハンドリングとリトライロジック
  - 詳細なリファレンス：`references/git_commands.md`、`references/best_practices.md`
