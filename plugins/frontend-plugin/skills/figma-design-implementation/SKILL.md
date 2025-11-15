---
name: figma-design-implementation
description: このスキルは、Figma デザインファイルから Vue コンポーネントの HTML/CSS を実装する際に使用するスキルです。Figma Dev Mode MCP を活用してデザイントークン（design token）を抽出し、プロジェクト規約に従ってコードを生成します。ユーザーが Figma デザインの実装、Vue コンポーネント（Vue component）の作成、またはデザインシステムのコード化を依頼した場合に使用してください。
---

# Figma デザイン実装スキル

## 概要

このスキルは、Figma デザインファイルから Vue プロジェクト向けの HTML/CSS コードを生成するためのワークフローとツールを提供します。Figma Dev Mode MCP を使用してデザイン情報を取得し、プロジェクト固有の規約に従って実装します。

**重要な原則:**
- **デザインの再現に集中** - ロジックの実装は行わず、Figma デザインの視覚的な忠実性を最優先
- **デザイントークンの活用** - CSS カスタムプロパティとして定義し、一貫性を保つ
- **プロジェクト規約の遵守** - `references/project-convention.md` に定義された規則に厳密に従う

---

## 実装ワークフロー

Figma デザインから Vue コンポーネントを実装する際は、以下の手順に従ってください。

### 1. デザイン情報の収集

#### 1.1 Figma ファイル/ノードの確認

ユーザーから Figma URL または選択中のデザインが提供されます。

**Figma URL の例:**
```
https://www.figma.com/file/ABC123XYZ/Design?node-id=1-2
```

#### 1.2 デザイントークンの抽出

MCP ツール `get_variable_defs` を使用してデザイントークンを取得します：

```python
# Claude が MCP ツールを呼び出す例
# mcp__figma_dev_mode_mcp__get_variable_defs()
```

**取得できる情報:**
- カラーパレット（primary, secondary, text など）
- スペーシング値（margins, padding など）
- タイポグラフィ（font-family, font-size, font-weight など）
- エフェクト（shadows, borders など）

#### 1.3 スクリーンショットの取得

視覚的なリファレンスとして、`get_screenshot` でデザインのスクリーンショットを取得します：

```python
# mcp__figma_dev_mode_mcp__get_screenshot()
```

これにより、細かいレイアウトやグラデーション、複雑なデザインの忠実な再現が可能になります。

#### 1.4 デザインコンテキストの取得

`get_design_context` でデザインを初期コードとして取得します：

```python
# Vue + CSS で取得する例
# mcp__figma_dev_mode_mcp__get_design_context({
#   "framework": "vue",
#   "styling": "css"
# })
```

**カスタマイズオプション:**
- `framework`: "vue", "react", "html"
- `styling`: "css", "tailwind", "scss"
- `custom_instructions`: 追加の指示（例: "BEM命名規則を使用"）

### 2. デザイントークンの変換

取得したデザイントークンを CSS カスタムプロパティに変換します。

**ツールの使用:**

`scripts/design_tokens_to_css.py` を使用：

```python
from scripts.design_tokens_to_css import design_tokens_to_css

# 取得したトークンを変換
tokens = {
    "colors": {"primary": "#3B82F6"},
    "spacing": {"md": "16px"}
}

css_variables = design_tokens_to_css(tokens)
```

**出力例:**
```css
:root {
  --color-primary: #3B82F6;
  --spacing-md: 16px;
}
```

これを `styles/tokens.css` として保存します。

### 3. コンポーネント構造の実装

#### 3.1 プロジェクト規約の確認

実装前にプロジェクトローカルの `${pwd}/.claude/skills/figma-design-implementation/references/project-convention.md` を読み込み、規約を確認してください：

- BEM 命名規則
- Vue SFC 構造
- CSS 変数の使用方法
- アクセシビリティ要件

**重要:** この規約は **必ず遵守** してください。

ファイルが存在しない場合は、`/fit-to-project figma-design-implementation` コマンドを実行してプロジェクト固有の skill を作成してください。

#### 3.2 Vue コンポーネントの作成

**テンプレート構造:**

```vue
<template>
  <div class="component-name">
    <div class="component-name__header">
      <!-- ヘッダー部分 -->
    </div>
    <div class="component-name__content">
      <!-- コンテンツ部分 -->
    </div>
  </div>
</template>

<script setup lang="ts">
// Props とロジックは後で実装（デザイン再現フェーズでは不要）
</script>

<style scoped>
/* BEM 命名規則とCSS変数を使用 */
.component-name {
  padding: var(--spacing-md);
  background-color: var(--color-background);
}

.component-name__header {
  margin-bottom: var(--spacing-lg);
  color: var(--color-text-primary);
}
</style>
```

#### 3.3 スタイリングのガイドライン

1. **CSS 変数を優先**
   - ハードコードされた値は使用しない
   - `styles/tokens.css` で定義された変数を参照

2. **BEM 命名規則を適用**
   ```css
   .block__element--modifier
   ```

3. **Scoped CSS を使用**
   - `<style scoped>` でコンポーネントスタイルを隔離

4. **レスポンシブデザイン**
   ```css
   .component {
     padding: var(--spacing-md);
   }

   @media (min-width: 768px) {
     .component {
       padding: var(--spacing-lg);
     }
   }
   ```

### 4. 検証とレビュー

#### 4.1 視覚的な一致の確認

取得したスクリーンショットと実装を比較し、以下を確認：
- レイアウトの正確性
- カラー、フォント、スペーシングの一致
- ボーダー、シャドウなどの細部

#### 4.2 チェックリストの確認

プロジェクトローカルの `${pwd}/.claude/skills/figma-design-implementation/references/project-convention.md` の最後にあるチェックリスト（存在する場合）を実行：

- [ ] BEM 命名規則に従っているか
- [ ] CSS 変数を使用しているか
- [ ] `<style scoped>` を使用しているか
- [ ] セマンティック HTML を使用しているか
- [ ] アクセシビリティ要件を満たしているか
- [ ] Figma デザインと視覚的に一致しているか

---

## Figma MCP ツールリファレンス

このスキルは `scripts/servers/figma/` ディレクトリに Figma Dev Mode MCP ツールのラッパーを提供しています。

### 利用可能なツール

各ツールの詳細は対応するファイルを参照してください：

1. **get_design_context** (`scripts/servers/figma/get_design_context.py`)
   - デザインをコードに変換
   - Vue, React, HTML 出力に対応

2. **get_variable_defs** (`scripts/servers/figma/get_variable_defs.py`)
   - デザイントークンを取得
   - 色、間隔、タイポグラフィなど

3. **get_screenshot** (`scripts/servers/figma/get_screenshot.py`)
   - デザインのスクリーンショットを取得
   - 視覚的リファレンスとして使用

4. **get_code_connect_map** (`scripts/servers/figma/get_code_connect_map.py`)
   - デザインとコードのマッピング
   - 既存コンポーネントの再利用

5. **create_design_system_rules** (`scripts/servers/figma/create_design_system_rules.py`)
   - デザインシステムルールを生成
   - プロジェクトの規約を理解

### ツールの探索

必要に応じて `scripts/servers/figma/` ディレクトリを探索し、各ツールファイルを読み込んで詳細を確認してください。

```python
# ツール一覧を確認
from scripts.servers.figma.index import list_available_tools, get_tool_info

tools = list_available_tools()
# ['get_design_context', 'get_variable_defs', ...]

info = get_tool_info('get_design_context')
# ツールの説明と用途を取得
```

---

## ユーティリティスクリプト

### design_tokens_to_css.py

デザイントークンを CSS に変換するユーティリティ。

**主要関数:**

1. **design_tokens_to_css()** - CSS カスタムプロパティに変換
2. **tokens_to_tailwind_config()** - Tailwind CSS 設定に変換
3. **tokens_to_scss_variables()** - SCSS 変数に変換

**使用例:**
```python
from scripts.design_tokens_to_css import design_tokens_to_css

tokens = {...}  # get_variable_defs から取得
css = design_tokens_to_css(tokens)

# styles/tokens.css に保存
with open('styles/tokens.css', 'w') as f:
    f.write(css)
```

---

## 実装例

### 例1: ボタンコンポーネント

**Figma デザイン:** プライマリボタン（青背景、白テキスト、角丸）

**実装手順:**

1. デザイントークンを取得し、CSS 変数として定義
2. BEM 命名規則でクラスを定義
3. Vue SFC を作成

**結果:**

```vue
<template>
  <button class="button button--primary">
    <span class="button__text">クリック</span>
  </button>
</template>

<script setup lang="ts">
// Props は後で実装
</script>

<style scoped>
.button {
  padding: var(--spacing-sm) var(--spacing-md);
  border: none;
  border-radius: var(--border-radius-md);
  font-family: var(--font-family-primary);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: background-color 0.2s;
}

.button--primary {
  background-color: var(--color-primary);
  color: var(--color-white);
}

.button--primary:hover {
  background-color: var(--color-primary-dark);
}

.button__text {
  display: inline-block;
}
</style>
```

### 例2: カードコンポーネント

**Figma デザイン:** ユーザープロフィールカード

**実装:**

```vue
<template>
  <div class="user-card">
    <div class="user-card__header">
      <img class="user-card__avatar" :src="avatarUrl" alt="プロフィール画像">
      <h3 class="user-card__name">{{ name }}</h3>
    </div>
    <div class="user-card__content">
      <p class="user-card__bio">{{ bio }}</p>
    </div>
  </div>
</template>

<style scoped>
.user-card {
  padding: var(--spacing-lg);
  background-color: var(--color-surface);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-md);
}

.user-card__header {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.user-card__avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
}

.user-card__name {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
  margin: 0;
}

.user-card__bio {
  font-size: var(--font-size-base);
  color: var(--color-text-secondary);
  line-height: 1.6;
}
</style>
```

---

## よくあるシナリオ

### シナリオ 1: 新しいページの実装

**ユーザーリクエスト:** "この Figma デザインを Vue コンポーネントとして実装してください"

**アプローチ:**
1. デザイントークンを抽出して `styles/tokens.css` を更新
2. スクリーンショットを取得
3. `get_design_context` で初期コードを生成
4. プロジェクト規約に従って調整
5. レイアウトとスタイルを実装
6. スクリーンショットと比較して検証

### シナリオ 2: 既存コンポーネントの更新

**ユーザーリクエスト:** "ボタンのデザインが更新されたので、コードも更新してください"

**アプローチ:**
1. `get_code_connect_map` で既存のマッピングを確認
2. 新しいデザイントークンを取得
3. 既存の CSS を更新（CSS 変数の値を変更）
4. 視覚的な一致を確認

### シナリオ 3: デザインシステムの構築

**ユーザーリクエスト:** "Figma のデザインシステムをコードとして実装したい"

**アプローチ:**
1. `create_design_system_rules` でデザインシステムルールを生成
2. すべてのデザイントークンを抽出
3. コンポーネントライブラリを実装
4. `references/design-system-rules.md` としてルールを保存

---

## 注意事項

### このフェーズで行うこと

- ✅ Figma デザインの視覚的な再現
- ✅ HTML/CSS の実装
- ✅ デザイントークンの抽出と定義
- ✅ レスポンシブレイアウト
- ✅ アクセシビリティ基本要件

### このフェーズで行わないこと

- ❌ ロジックの実装（JavaScript/TypeScript コード）
- ❌ API 統合
- ❌ 状態管理
- ❌ イベントハンドラーの実装（後で追加）

**理由:** デザイン再現フェーズでは、視覚的な忠実性に集中し、機能的な実装は別のフェーズで行います。

---

## リソース

### プロジェクトローカル: ${pwd}/.claude/skills/figma-design-implementation/

- **references/project-convention.md** - プロジェクト固有の規約（必読、低自由度）
  - BEM 命名規則
  - Vue コンポーネント構造
  - CSS スタイリング規約
  - アクセシビリティ要件
  - その他のプロジェクト固有のリソース

### マーケットプレイス: このスキル

- **scripts/** - Figma MCP ツールラッパーとユーティリティ
  - `servers/figma/` - Figma MCP ツールラッパー
    - `get_design_context.py` - デザインからコード生成
    - `get_variable_defs.py` - デザイントークン取得
    - `get_screenshot.py` - スクリーンショット取得
    - その他のツール
  - `design_tokens_to_css.py` - トークン変換ユーティリティ

- **references/project-convention.md** - テンプレート（プロジェクト固有の skill 作成時に使用）

- **assets/** - 現在は空だが、将来的にボイラープレートやテンプレートを追加可能

プロジェクトローカルの skill が作成されている場合は、そちらの references を優先的に参照する。

---

## トラブルシューティング

### MCP ツールが呼び出せない

- Figma Dev Mode MCP サーバーが起動しているか確認
- `.mcp.json` で正しく設定されているか確認（SSE, http://127.0.0.1:3845/sse）

### デザインが正確に再現できない

- `get_screenshot` でスクリーンショットを取得し、視覚的リファレンスとして使用
- 複雑なグラデーションやエフェクトは CSS で近似

### CSS 変数が定義されていない

- `get_variable_defs` でデザイントークンを再取得
- `design_tokens_to_css.py` で CSS 変数を生成
- `styles/tokens.css` に保存されているか確認
