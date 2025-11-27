# 推奨コマンド

## プロジェクト管理

### Git操作
```bash
# ブランチ作成と切り替え
git switch -c <branch-name>

# 変更をステージングしてコミット
git add .
git commit -m "commit message"

# リモートにプッシュ
git push -u origin <branch-name>

# ブランチ一覧確認
git branch -a

# 現在のブランチ確認
git branch --show-current
```

### ステータス確認
```bash
# Gitステータス確認
git status

# 変更内容の確認
git diff

# コミット履歴確認
git log --oneline
```

## library-docs-reference スキル

### セットアップ
```bash
# スキルディレクトリに移動
cd plugins/development-plugin/skills/library-docs-reference/scripts

# 依存関係のインストール
npm install
```

### ドキュメント取得
```bash
# 基本的な使い方
tsx context7-workflow.ts <library-name> [topic] [mode] [page]

# 例: Reactのドキュメントを取得
tsx context7-workflow.ts react

# 例: ReactのHooksドキュメントを取得
tsx context7-workflow.ts react hooks

# 例: Next.jsのルーティングガイドを取得（情報モード）
tsx context7-workflow.ts next.js routing info 1
```

### npm スクリプト
```bash
# 基本的な使用例を実行
npm run example:basic

# 複数ページのドキュメントを取得
npm run example:multi

# ライブラリを比較
npm run example:compare react vue --topic hooks
```

## プロジェクト探索

### ディレクトリ構造確認
```bash
# プロジェクトルートの構造を確認
ls -la

# プラグイン一覧を確認
ls -la plugins/

# スキル一覧を確認
ls -la plugins/development-plugin/skills/
ls -la plugins/frontend-plugin/skills/
```

### ファイル検索
```bash
# 特定のファイルを検索
find . -name "*.json" -not -path "*/node_modules/*"
find . -name "SKILL.md"

# 特定の文字列を検索
grep -r "keyword" --exclude-dir=node_modules --exclude-dir=.git
```

## システムユーティリティ（Linux）

### ファイル操作
```bash
# ファイル内容を表示
cat <file>

# ファイルの先頭を表示
head -n 20 <file>

# ファイルの末尾を表示
tail -n 20 <file>

# ディレクトリのサイズを確認
du -sh <directory>
```

### プロセス管理
```bash
# 実行中のプロセスを確認
ps aux | grep <process-name>

# ポートの使用状況を確認
lsof -i :<port-number>
```

## Serena MCP関連

### プロジェクトアクティベーション
プロジェクトをSerenaに認識させる場合、MCP toolを使用してアクティベート。

### メモリー管理
プロジェクト情報はメモリーファイル（`.serena/memories/`）に保存され、将来のタスクで参照可能。

## 注意事項

1. **言語**: すべてのやり取りは日本語で行う
2. **MCP設定**: Context7のAPIキーは `.mcp.json` で管理（セキュリティに注意）
3. **依存関係**: TypeScriptスキルを使う前に `npm install` を実行
4. **リソース管理**: MCP接続は使用後に必ずクローズ
