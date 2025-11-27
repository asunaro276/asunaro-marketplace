# プロジェクト概要

## プロジェクト名
asunaro-marketplace

## プロジェクトの目的
Claude Codeで使用するプラグインのマーケットプレイス。複数のプラグインとスキルを管理し、開発作業を効率化するためのリポジトリ。

## プロジェクト構造

```
asunaro-marketplace/
├── plugins/
│   ├── development-plugin/     # 開発全般に必要なプラグイン
│   │   ├── skills/
│   │   │   ├── git-operations/              # Git操作ガイド
│   │   │   └── library-docs-reference/      # Context7によるライブラリドキュメント取得
│   │   └── .claude-plugin/
│   │       └── plugin.json
│   │
│   └── frontend-plugin/        # フロントエンド開発用プラグイン
│       ├── skills/
│       │   ├── component-tdd/               # コンポーネントTDD開発
│       │   └── figma-design-implementation/ # Figmaデザイン実装
│       ├── commands/
│       └── .claude-plugin/
│           └── plugin.json
│
├── CLAUDE.md                    # Claude向けプロジェクト指示
└── README.md                    # (未作成)

```

## プラグイン一覧

### 1. development-plugin
**説明**: 開発を行う際に言語にかかわらず必要となるClaude Plugin

**スキル**:
- `git-operations`: Git操作に関する包括的なガイド（ブランチ、コミット、プッシュ、プル、マージ、リベースなど）
- `library-docs-reference`: Context7 MCPを使用してライブラリのドキュメントとコード例を取得

### 2. frontend-plugin
**説明**: フロントエンド開発を行う際に必要となる Claude Plugin。Figma デザインから Vue/HTML/CSS への実装をサポート

**スキル**:
- `component-tdd`: Vue、React、AstroなどのフロントエンドコンポーネントをTDDで作成
- `figma-design-implementation`: Figma デザインファイルから Vue コンポーネントの HTML/CSS を実装

## 技術スタック

### プラグインシステム
- Claude Code Plugin System
- MCP (Model Context Protocol)

### 使用しているMCPサーバー

**development-plugin:**
- Serena MCP (コードベース分析とシンボリック操作)
- Context7 MCP (ライブラリドキュメント取得)

**frontend-plugin:**
- Serena MCP (コードベース分析とシンボリック操作)
- Figma Dev Mode MCP (Figmaデザイン情報取得)

### library-docs-reference スキルの技術スタック
- TypeScript
- Node.js
- @modelcontextprotocol/sdk
- tsx (TypeScript実行環境)
