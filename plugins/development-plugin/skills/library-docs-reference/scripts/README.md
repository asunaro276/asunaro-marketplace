# Context7 Workflow Scripts

Context7 MCPサーバーを使用してライブラリのドキュメントを取得するためのワークフロースクリプト集です。

## ファイル構成

- `mcp-client.ts` - MCP クライアントの基本実装
- `context7-workflow.ts` - Context7を使用したドキュメント取得ワークフロー
- `examples/` - 使用例スクリプト

## セットアップ

### 前提条件

- Node.js 18以上
- TypeScript
- `@modelcontextprotocol/sdk` パッケージ

### インストール

```bash
npm install @modelcontextprotocol/sdk
```

## 使い方

### 1. TypeScriptを直接実行

```bash
# 基本的な使用方法
tsx context7-workflow.ts react

# トピックを指定
tsx context7-workflow.ts react hooks

# モードとページを指定
tsx context7-workflow.ts react hooks code 1

# 情報モードでドキュメントを取得
tsx context7-workflow.ts next.js routing info 1
```

### 2. JavaScriptにビルドして実行

```bash
# TypeScriptをビルド
tsc context7-workflow.ts --module esnext --target es2020 --moduleResolution node

# JavaScriptを実行
node context7-workflow.js vue
```

### 3. プログラムから使用

```typescript
import {
  resolveLibraryId,
  getLibraryDocs,
  getDocsByLibraryName,
  getMultiplePages
} from './context7-workflow.js';

// ライブラリIDを解決
const resolved = await resolveLibraryId('react');
console.log(resolved.selected?.id); // '/facebook/react'

// ドキュメントを取得
const docs = await getLibraryDocs('/facebook/react', 'hooks', 'code', 1);
console.log(docs.content);

// ライブラリ名から直接ドキュメントを取得（推奨）
const docs2 = await getDocsByLibraryName('vue', 'composition-api', 'code', 1);
console.log(docs2.content);

// 複数ページを取得
const allDocs = await getMultiplePages('/vercel/next.js', 'routing', 'code', 3);
allDocs.forEach((doc, i) => {
  console.log(`Page ${i + 1}:`, doc.content);
});
```

## API リファレンス

### `resolveLibraryId(libraryName: string)`

ライブラリ名からContext7互換のライブラリIDを解決します。

**パラメータ:**
- `libraryName` - 検索するライブラリ名

**戻り値:**
- `matches` - マッチしたライブラリのリスト
- `selected` - 選択されたライブラリ（最も関連性が高い）

### `getLibraryDocs(libraryId, topic?, mode?, page?)`

ライブラリIDからドキュメントを取得します。

**パラメータ:**
- `libraryId` - Context7互換のライブラリID (例: '/mongodb/docs', '/vercel/next.js')
- `topic` - オプション: ドキュメントのトピック (例: 'hooks', 'routing')
- `mode` - 'code' (デフォルト) または 'info'
  - `code`: APIリファレンスとコード例
  - `info`: 概念的なガイドとアーキテクチャ情報
- `page` - ページ番号 (1から開始、最大10)

**戻り値:**
- `content` - ドキュメントの内容
- `page` - 現在のページ番号
- `hasMore` - さらにページがあるかどうか

### `getDocsByLibraryName(libraryName, topic?, mode?, page?)`

ライブラリ名からドキュメントを取得する完全なワークフロー。
内部で `resolveLibraryId` を呼び出してからドキュメントを取得します。

**パラメータ:**
- `libraryName` - ライブラリ名
- `topic` - オプション: トピック
- `mode` - 'code' または 'info'
- `page` - ページ番号

### `getMultiplePages(libraryId, topic, mode?, maxPages?)`

複数ページのドキュメントを一度に取得します。

**パラメータ:**
- `libraryId` - ライブラリID
- `topic` - トピック
- `mode` - 'code' または 'info'
- `maxPages` - 取得する最大ページ数 (最大10)

## 使用例

### React Hooksのドキュメントを取得

```bash
tsx context7-workflow.ts react hooks
```

### Next.js のルーティングガイドを取得

```bash
tsx context7-workflow.ts next.js routing info
```

### Vue.js のComposition APIドキュメントを複数ページ取得

```typescript
import { getMultiplePages, closeConnection } from './context7-workflow.js';

const docs = await getMultiplePages('/vuejs/vue', 'composition-api', 'code', 3);
docs.forEach((doc, i) => console.log(`=== Page ${i + 1} ===\n${doc.content}`));
await closeConnection();
```

## トラブルシューティング

### 接続エラー

`.mcp.json` の `context7` 設定を確認してください：

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

### ライブラリが見つからない

ライブラリ名を変えて試してください：
- "react" → "React" または "facebook/react"
- "nextjs" → "next.js" または "next"
- "vue" → "Vue" または "vuejs"

## 参考リンク

- [Context7 Documentation](https://context7.com)
- [MCP SDK](https://github.com/modelcontextprotocol/sdk)
