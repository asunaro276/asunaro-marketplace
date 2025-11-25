# Context7 MCPサーバー ツールリファレンス

このドキュメントは、Context7 MCPサーバーで利用可能なツールの詳細な仕様とパラメータを提供する。

## 概要

Context7は、LLMとAIコードエディタに最新のコードドキュメントを提供するMCPサーバー。バージョン固有のドキュメントとコード例をソースから直接取得し、プロンプトに配置することで、幻覚APIを防ぐ。

## 利用可能なツール

### 1. resolve-library-id

ライブラリ名をContext7互換のライブラリIDに解決し、一致するライブラリのリストを返す。

#### 目的

`get-library-docs`を呼び出す前に、有効なContext7互換のライブラリIDを取得するために使用する。ただし、ユーザーが明示的に`/org/project`または`/org/project/version`形式のライブラリIDを提供している場合は、このツールを呼び出す必要はない。

#### パラメータ

| パラメータ名 | 型 | 必須 | 説明 |
|------------|------|------|------|
| `libraryName` | string | ✓ | 検索するライブラリ名（例: "react", "next.js", "mongodb"） |

#### 選択プロセス

1. クエリを分析して、ユーザーが探しているライブラリを理解する
2. 以下の基準に基づいて、最も関連性の高い一致を返す：
   - **名前の類似性**: クエリとの名前の一致度（完全一致を優先）
   - **説明の関連性**: クエリの意図に対する説明の関連性
   - **ドキュメントのカバレッジ**: コードスニペット数が多いライブラリを優先
   - **ソースの評判**: High または Medium の評判を持つライブラリを優先
   - **ベンチマークスコア**: 品質指標（100が最高スコア）

#### レスポンス形式

- 選択されたライブラリIDを明確にマークされたセクションで返す
- このライブラリが選択された理由について簡潔な説明を提供する
- 複数の良い一致が存在する場合は、それを認識するが、最も関連性の高いものを選択する
- 良い一致が存在しない場合は、それを明確に述べ、クエリの改善を提案する

#### 使用例

```python
# Reactライブラリの解決
result = resolve_library_id(libraryName="react")
# 返されるID例: "/facebook/react"

# Next.jsライブラリの解決
result = resolve_library_id(libraryName="next.js")
# 返されるID例: "/vercel/next.js"

# MongoDBドキュメントの解決
result = resolve_library_id(libraryName="mongodb")
# 返されるID例: "/mongodb/docs"
```

#### エラーハンドリング

- **曖昧なクエリ**: 複数の候補がある場合、最も関連性の高いものを選択するが、他の候補も言及する
- **一致なし**: ライブラリが見つからない場合、クエリの改善方法を提案する
- **名前の変動**: 一般的な名前のバリエーション（例: "nextjs" vs "next.js"）を処理する

---

### 2. get-library-docs

指定されたライブラリの最新ドキュメントを取得する。`resolve-library-id`で取得したContext7互換のライブラリIDを使用する必要がある（ユーザーが明示的にIDを提供している場合を除く）。

#### 目的

特定のライブラリのAPIリファレンス、コード例、概念的ガイド、アーキテクチャ情報を取得する。

#### パラメータ

| パラメータ名 | 型 | 必須 | デフォルト | 説明 |
|------------|------|------|-----------|------|
| `context7CompatibleLibraryID` | string | ✓ | - | `resolve-library-id`から取得した、またはユーザーが提供したContext7互換のライブラリID（例: `/mongodb/docs`, `/vercel/next.js`, `/vercel/next.js/v14.3.0-canary.87`） |
| `mode` | enum | × | `"code"` | ドキュメントモード：<br>- `"code"`: APIリファレンスとコード例用<br>- `"info"`: 概念的ガイド、アーキテクチャ情報用 |
| `topic` | string | × | - | フォーカスするトピック（例: `"hooks"`, `"routing"`, `"authentication"`） |
| `page` | integer | × | 1 | ページネーション用のページ番号（1から開始）。コンテキストが不十分な場合は`page=2`, `page=3`, `page=4`などを試す。最大10ページまで。 |

#### モードの選択

- **`mode='code'`（デフォルト）**:
  - APIリファレンス
  - 関数・メソッドのシグネチャ
  - コード例とスニペット
  - 使用パターン

- **`mode='info'`**:
  - 概念的なガイド
  - アーキテクチャの説明
  - ベストプラクティス
  - 設計哲学

#### トピックの指定

効果的なトピック指定の例：

**具体的な機能**:
- `"hooks"` - Reactのフック
- `"routing"` - Next.jsのルーティング
- `"authentication"` - 認証機能
- `"middleware"` - ミドルウェア

**コンポーネント名**:
- `"Button"` - ボタンコンポーネント
- `"Modal"` - モーダルコンポーネント
- `"Form"` - フォームコンポーネント

**概念**:
- `"state management"` - 状態管理
- `"data fetching"` - データ取得
- `"error handling"` - エラーハンドリング

#### 使用例

```python
# 例1: ReactのフックのAPIリファレンスとコード例を取得
docs = get_library_docs(
    context7CompatibleLibraryID="/facebook/react",
    mode="code",
    topic="hooks",
    page=1
)

# 例2: Next.jsのルーティングに関する概念的ガイドを取得
docs = get_library_docs(
    context7CompatibleLibraryID="/vercel/next.js",
    mode="info",
    topic="routing",
    page=1
)

# 例3: MongoDBのクエリ構文のコード例を取得
docs = get_library_docs(
    context7CompatibleLibraryID="/mongodb/docs",
    mode="code",
    topic="query syntax",
    page=1
)

# 例4: 特定バージョンのNext.jsドキュメントを取得
docs = get_library_docs(
    context7CompatibleLibraryID="/vercel/next.js/v14.3.0-canary.87",
    mode="code",
    topic="app router"
)

# 例5: コンテキストが不十分な場合、追加ページを取得
docs_page2 = get_library_docs(
    context7CompatibleLibraryID="/facebook/react",
    mode="code",
    topic="hooks",
    page=2
)
```

#### ページネーション

- `page`パラメータを使用して、複数ページのドキュメントを順次取得できる
- コンテキストが不十分な場合は、同じトピックで`page=2`, `page=3`などを試す
- 最大10ページまで取得可能（`page=1`から`page=10`）
- 各ページには関連するドキュメントの異なるセクションが含まれる

#### エラーハンドリング

- **無効なライブラリID**: ライブラリIDが見つからない場合、`resolve-library-id`を再度実行して確認する
- **トピックが見つからない**: トピックの指定を変更するか、より一般的なトピックを試す
- **ページ範囲外**: 有効なページ番号（1-10）を指定する

---

## 統合ワークフロー

### 基本的な使用パターン

```python
# ステップ1: ライブラリIDを解決（ユーザーがIDを提供していない場合）
library_id_result = resolve_library_id(libraryName="react")
library_id = library_id_result["id"]  # 例: "/facebook/react"

# ステップ2: ドキュメントを取得
docs = get_library_docs(
    context7CompatibleLibraryID=library_id,
    mode="code",
    topic="hooks"
)

# ステップ3: 取得したドキュメントを使用してコードを実装
# （ドキュメントの内容をプロジェクトに適用）
```

### ユーザーがライブラリIDを提供している場合

```python
# ユーザーが "/vercel/next.js" を提供している場合
# resolve-library-id をスキップして直接 get-library-docs を使用
docs = get_library_docs(
    context7CompatibleLibraryID="/vercel/next.js",
    mode="code",
    topic="routing"
)
```

### 複数トピックの処理

```python
# 複数のトピックについてドキュメントを取得
topics = ["hooks", "context", "state management"]

for topic in topics:
    docs = get_library_docs(
        context7CompatibleLibraryID="/facebook/react",
        mode="code",
        topic=topic
    )
    # 各トピックのドキュメントを処理
```

---

## ベストプラクティス

### 1. ライブラリIDの解決

- ユーザーが明示的にライブラリIDを提供していない場合は、必ず`resolve-library-id`を最初に呼び出す
- 曖昧なライブラリ名の場合は、解決結果を確認してから次に進む

### 2. モードの選択

- APIの使い方やコード例が必要な場合は`mode='code'`を使用
- 概念的な理解やアーキテクチャが必要な場合は`mode='info'`を使用
- 両方が必要な場合は、2回呼び出す

### 3. トピックの絞り込み

- 具体的なトピックを指定して、関連するドキュメントのみを取得
- トピックは、ライブラリの一般的な用語や機能名を使用
- 結果が不十分な場合は、トピックを変更または一般化する

### 4. ページネーションの活用

- 最初のページで十分な情報が得られない場合は、追加ページを取得
- 同じトピックで複数ページを順次取得して、より包括的な情報を得る

### 5. コンテキスト効率

- 大規模なドキュメントを取得する際は、トピックを絞り込む
- 不要なドキュメントを取得しないように、クエリを具体的にする
- ページネーションを使用して、必要な情報を段階的に取得

---

## トラブルシューティング

### 問題: ライブラリが見つからない

**解決策**:
- ライブラリ名のスペルを確認
- 一般的な名前のバリエーションを試す（例: "nextjs" → "next.js"）
- ライブラリの公式名称を使用

### 問題: トピックに関するドキュメントが見つからない

**解決策**:
- トピック名を変更（例: "state" → "state management"）
- より一般的なトピックを試す
- `mode`を変更（`code` ⇔ `info`）
- ページネーションを使用して追加ページを確認

### 問題: 取得したドキュメントが不十分

**解決策**:
- `page`パラメータを使用して追加ページを取得
- トピックを変更してより関連性の高い情報を検索
- 複数のトピックで個別にドキュメントを取得

### 問題: 古いバージョンのドキュメントが返される

**解決策**:
- ライブラリIDにバージョンを明示的に指定（例: `/vercel/next.js/v14.3.0-canary.87`）
- `resolve-library-id`で最新バージョンを確認

---

## 参考リンク

- [Context7 GitHub Repository](https://github.com/upstash/context7)
- [MCP Code Execution Best Practices](https://www.anthropic.com/engineering/code-execution-with-mcp)
