---
name: library-docs-reference
description:
   context7によるライブラリ情報の取得を完全にコード内で実行。
   MCP Toolsの処理回数を削減
   ユーザーが言及したライブラリのドキュメントやコード例を参照したい時に使用。
   公式で言及されている最新のベストプラクティスを発見。
---

# **禁止事項**
context7 mcpの直接実行は決してしないでください
必ずスクリプトの実行を通してMCPを利用してください

# Context7ワークフローの実行手順

このスキルは、Context7 MCPサーバーを使用してライブラリのドキュメントを取得します。
`scripts/` ディレクトリ配下のワークフロースクリプトを使用して、効率的にドキュメントを取得できます。

## セットアップ

### 1. 依存関係のインストール

```bash
cd plugins/development-plugin/skills/library-docs-reference/scripts
npm install
```

必要なパッケージ：
- `@modelcontextprotocol/sdk` - MCP SDK
- `@types/node` - Node.js型定義
- `tsx` - TypeScript実行環境
- `typescript` - TypeScriptコンパイラ

### 2. MCP設定の確認

`plugins/development-plugin/.mcp.json` に以下の設定があることを確認：

```json
{
  "mcpServers": {
    "context7": {
      "type": "http",
      "url": "https://mcp.context7.com/mcp",
      "headers": {
        "CONTEXT7_API_KEY": "your-api-key"
      }
    }
  }
}
```

## ワークフローの実行方法

### 方法1: npmスクリプト実行

```bash
cd plugins/development-plugin/skills/library-docs-reference/scripts
```

### 方法2: プログラムから使用

```typescript
import {
  resolveLibraryId,
  getLibraryDocs,
  getDocsByLibraryName,
  closeConnection
} from './scripts/context7-workflow.js';

// ライブラリ名から直接ドキュメントを取得（推奨）
const docs = await getDocsByLibraryName('react', 'hooks', 'code', 1);
console.log(docs.content);

// 接続をクローズ
await closeConnection();
```

## 主な関数

### `resolveLibraryId(libraryName: string)`

ライブラリ名からContext7互換のライブラリIDを解決します。

```typescript
const result = await resolveLibraryId('react');
// { matches: [...], selected: { id: '/facebook/react', ... } }
```

### `getLibraryDocs(libraryId, topic?, mode?, page?)`

ライブラリIDからドキュメントを取得します。

パラメータ：
- `libraryId`: Context7互換のライブラリID (例: '/mongodb/docs', '/vercel/next.js')
- `topic`: ドキュメントのトピック (例: 'hooks', 'routing')
- `mode`: 'code'（デフォルト）または 'info'
  - `code`: APIリファレンスとコード例
  - `info`: 概念的なガイドとアーキテクチャ情報
- `page`: ページ番号（1から開始、最大10）

```typescript
const docs = await getLibraryDocs('/facebook/react', 'hooks', 'code', 1);
```

### `getDocsByLibraryName(libraryName, topic?, mode?, page?)`

ライブラリ名から直接ドキュメントを取得する完全なワークフロー（推奨）。

```typescript
const docs = await getDocsByLibraryName('vue', 'composition-api', 'code', 1);
```

### `getMultiplePages(libraryId, topic, mode?, maxPages?)`

複数ページのドキュメントを一度に取得します。

```typescript
const allDocs = await getMultiplePages(
  '/vercel/next.js',
  'routing',
  'code',
  3
);
```

## 使用例

### 例1: Reactのフックドキュメントを取得

```bash
tsx context7-workflow.ts react hooks
```

### 例4: プログラムから複数ページを処理

```typescript
import { getMultiplePages, closeConnection } from './scripts/context7-workflow.js';

try {
  const pages = await getMultiplePages('/facebook/react', 'hooks', 'code', 5);

  // 各ページを処理
  for (const [index, page] of pages.entries()) {
    console.log(`=== Page ${index + 1} ===`);
    console.log(page.content);
  }
} finally {
  await closeConnection();
}
```

## ワークフローの内部動作

### `getDocsByLibraryName` の処理フロー

1. **ライブラリID解決**: `resolveLibraryId()` を呼び出してライブラリ名からIDを取得
   - Context7 MCP の `resolve-library-id` ツールを使用
   - 複数の候補から最も関連性の高いものを選択

2. **ドキュメント取得**: `getLibraryDocs()` で実際のドキュメントを取得
   - Context7 MCP の `get-library-docs` ツールを使用
   - トピック、モード、ページを指定

3. **結果の返却**: ドキュメント内容とメタデータを返す

### MCPクライアントの仕組み

`mcp-client.ts` は以下の機能を提供：

1. **設定読み込み**: `.mcp.json` からContext7の設定を読み込み
2. **接続管理**: SSE（Server-Sent Events）を使用してHTTP接続を確立
3. **ツール呼び出し**: MCPツールを呼び出してレスポンスをパース
4. **接続キャッシュ**: 複数回の呼び出しで接続を再利用

```typescript
// mcp-client.ts の主要関数
export async function callMCPTool<T>(toolName: string, input: any): Promise<T>
export async function closeConnection(): Promise<void>
```

## トラブルシューティング

### エラー: "MCP server 'context7' not found"

`.mcp.json` の設定を確認してください。ファイルは `plugins/development-plugin/.mcp.json` に配置されている必要があります。

### エラー: "Cannot find name 'console'" または型エラー

```bash
cd scripts
npm install @types/node
```

### エラー: "No library found for: xxx"

ライブラリ名を変更して試してください：
- "react" → "React" または "facebook/react"
- "nextjs" → "next.js" または "next"
- "vue" → "Vue" または "vuejs"

### 接続がタイムアウトする

Context7 APIキーが正しいか確認してください。また、ネットワーク接続を確認してください。

## 詳細情報

より詳しい使い方については、以下を参照してください：
- `scripts/README.md` - 詳細なドキュメント
- [Context7 Documentation](https://context7.com)
