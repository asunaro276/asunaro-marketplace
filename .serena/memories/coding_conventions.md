# コーディング規約とスタイルガイド

## 全般的な規約

### 言語
- ユーザーとのやり取りはすべて日本語で行う（CLAUDE.mdに明記）

### ファイル構造
- プラグインは `plugins/` ディレクトリに配置
- 各プラグインには `.claude-plugin/plugin.json` が必要
- スキルは各プラグインの `skills/` ディレクトリに配置
- スキルには `SKILL.md` ファイルでドキュメント化

## プラグイン設定ファイル (plugin.json)

### 必須フィールド
```json
{
  "name": "plugin-name",
  "description": "プラグインの説明（日本語）",
  "version": "1.0.0",
  "author": {
    "name": "作者名",
    "url": "GitHub URL"
  },
  "mcpServers": {
    // MCP サーバー設定
  }
}
```

## スキル設定ファイル (SKILL.md)

### フロントマター形式
```markdown
---
name: skill-name
description: スキルの説明（日本語）
---
```

### 構成要素
1. 概要セクション
2. いつ使用するか
3. ワークフロー/実行手順
4. 例とサンプルコード
5. トラブルシューティング

## MCP 設定

### development-plugin の .mcp.json
- Serena MCP: コードベース分析用
- Context7 MCP: ライブラリドキュメント取得用（HTTPベース）

### frontend-plugin の plugin.json
- Serena MCP: コードベース分析用
- Figma Dev Mode MCP: Figmaデザイン取得用（SSEベース）

## TypeScript スキル開発（library-docs-reference）

### ファイル構成
```
scripts/
├── package.json
├── tsconfig.json
├── context7-workflow.ts    # メインワークフロー
├── mcp-client.ts           # MCPクライアント
└── examples/               # 使用例
```

### TypeScript設定
- ES Modules使用 (`"type": "module"`)
- Node.js型定義を含める (`@types/node`)
- tsx でTypeScriptを直接実行

### コーディングスタイル
- 関数はexportして再利用可能にする
- エラーハンドリングを適切に実装
- 接続リソースは必ずクローズする
- ドキュメントコメントを含める

## ドキュメント規約

### README.md / SKILL.md
- 日本語で記述
- セットアップ手順を明記
- 使用例を複数提供
- トラブルシューティングセクションを含める
- コマンドラインとプログラムの両方の使い方を示す

### コード例
- 実際に動作する例を提供
- コメントは日本語で記述
- 複数のユースケースをカバー
