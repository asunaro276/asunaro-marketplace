---
description: プロジェクト固有の慣習やアーキテクチャ、コーディング規約を分析して、汎用スキルにプロジェクト固有の情報を適用する
args:
  - name: skill_name
    description: カスタマイズ対象のスキル名（例：git-operations）
    required: true
---

# スキルのカスタマイズ

指定されたスキルに、プロジェクト固有の慣習、アーキテクチャ、コーディング規約を適用します。

## 処理フロー

### 1. プロジェクトドキュメントの分析
project-docs-analyzerエージェントを使用して、以下を抽出：
- 開発規約
- アーキテクチャパターン
- コーディング標準
- ツール構成
- 命名規則

対象ドキュメント：CLAUDE.md、README.md、その他のMarkdownファイル

### 2. コードベースの分析
codebase-analyzerエージェントを使用して、以下を抽出：
- 技術スタック
- ディレクトリ構造
- 実装パターン
- テスト手法

### 3. スキルへの適用
分析結果を統合し、`references/project-convention.md` を作成・更新

### 4. スキルの更新
skill-creatorを使用してスキル本体を更新し、プロジェクト固有のツールやベストプラクティスをreference配下にまとめる

### 5. クリーンアップ
中間ファイルを削除

---

# 実装手順

## ステップ1: 一時ディレクトリの作成

まず、分析結果を保存する一時ディレクトリを作成します。

```bash
mkdir -p .claude/tmp
```

## ステップ2: 並列分析の実行

Task toolを使用して、以下の2つのエージェントを**並列で**起動します：

### エージェント1: プロジェクトドキュメント分析

```
subagent_type: project-docs-analyzer
prompt: このプロジェクトのドキュメント（CLAUDE.md、README.md、その他のMarkdownファイル）を分析し、以下の情報を抽出してください：
- プロジェクト固有の開発規約
- アーキテクチャパターン
- コーディング標準
- ツール構成
- 命名規則

分析結果を構造化されたMarkdown形式で `.claude/tmp/project-docs-analysis.md` に保存してください。

出力形式：
# プロジェクトドキュメント分析結果

## プロジェクト概要
[概要]

## 開発規約
[開発規約の詳細]

## アーキテクチャパターン
[アーキテクチャの詳細]

## コーディング標準
[コーディング標準の詳細]

## ツール構成
[使用するツールと設定]

## 命名規則
[命名規則の詳細]

## その他
[その他の重要な情報]
```

### エージェント2: コードベース分析

```
subagent_type: codebase-analyzer
prompt: このプロジェクトのコードベースを分析し、以下の情報を抽出してください：
- 技術スタック（言語、フレームワーク、ライブラリ）
- ディレクトリ構造とその意図
- 実装パターン（設計パターン、アーキテクチャパターン）
- テスト手法とテストツール

分析結果を構造化されたMarkdown形式で `.claude/tmp/codebase-analysis.md` に保存してください。

出力形式：
# コードベース分析結果

## 技術スタック
[技術スタックの詳細]

## ディレクトリ構造
[ディレクトリ構造の詳細]

## 実装パターン
[実装パターンの詳細]

## テスト手法
[テスト手法の詳細]

## その他
[その他の重要な情報]
```

**重要**: 2つのエージェントは並列で実行してください。

## ステップ3: スキルの所属プラグインを特定

分析が完了したら、指定されたスキル（{skill_name}）がどのプラグインに属しているか特定します。

```bash
# スキルディレクトリを検索
find plugins -type d -name "{skill_name}" | grep "/skills/"
```

例：git-operationsスキルの場合
→ `plugins/development-plugin/skills/git-operations`
→ プラグイン名: development-plugin

## ステップ4: 分析結果の統合

`.claude/tmp/project-docs-analysis.md` と `.claude/tmp/codebase-analysis.md` を読み込み、{skill_name}スキルに関連する情報を抽出・統合します。

統合時の考慮事項：
- スキルのドメインに関連する情報を優先
- 汎用的な情報と具体的な情報をバランスよく含める
- プロジェクト固有のベストプラクティスを明確に記載
- 実装例やコマンド例を含める

## ステップ5: project-convention.mdの作成・更新

以下の内容で `plugins/{plugin_name}/skills/{skill_name}/references/project-convention.md` を作成・更新します：

```markdown
# プロジェクト固有の規約と慣習

このファイルは、{skill_name}スキルに適用されるプロジェクト固有の情報を含みます。

最終更新日: [YYYY-MM-DD]

---

## プロジェクト概要

[プロジェクトの概要と目的]

## 開発規約

### コーディング標準
[プロジェクト固有のコーディング標準]

### 命名規則
[プロジェクト固有の命名規則]

### ファイル構成
[プロジェクト固有のファイル構成ルール]

## アーキテクチャ

### 技術スタック
[使用している技術スタック]

### アーキテクチャパターン
[採用しているアーキテクチャパターン]

### ディレクトリ構造
[プロジェクトのディレクトリ構造と意図]

## {skill_name}に関連するプロジェクト固有の情報

### ベストプラクティス
[このスキルに関連するプロジェクト固有のベストプラクティス]

### 注意事項
[このスキルを使用する際の特別な注意事項]

### 実装例
[プロジェクトでの実装例]

## ツールと設定

### 使用ツール
[プロジェクトで使用するツール]

### 設定ファイル
[関連する設定ファイルと内容]

## テスト

### テスト手法
[プロジェクトで採用しているテスト手法]

### テストツール
[使用しているテストツール]

## 参考資料

- プロジェクトドキュメント: [リンクまたはパス]
- 関連する設定ファイル: [パス]
```

**重要**: 既存の `project-convention.md` がある場合は上書きされます。

## ステップ6: スキルの更新（skill-creatorを使用）

作成した `project-convention.md` の内容を基に、スキル本体を更新します。

### 6.1 skill-creatorスキルの使用

メインエージェントがスキルの書き換えを行う際は、必ず**skill-creator**スキルを使用してください。

```
Skill tool を使用:
skill: "skill-creator"

プロンプト:
{skill_name}スキルを以下の方針で更新してください：

1. project-convention.mdの内容を参照し、プロジェクト固有の情報をスキルに統合
2. スキルの説明や使用例をプロジェクトに適したものに調整
3. 汎用性を保ちつつ、プロジェクト固有のベストプラクティスを反映

対象スキル: plugins/{plugin_name}/skills/{skill_name}/SKILL.md
参照ファイル: plugins/{plugin_name}/skills/{skill_name}/references/project-convention.md
```

### 6.2 reference配下へのベストプラクティスのまとめ

プロジェクトで利用されているツールやフレームワークのベストプラクティスについては、`references/` ディレクトリ配下に個別ファイルとしてまとめます。

**作成すべきreferenceファイルの例:**

1. **ツール固有のベストプラクティス**
   - `references/git-workflow.md` - プロジェクトのGitワークフロー
   - `references/testing-practices.md` - テストの書き方とツールの使用方法
   - `references/ci-cd-practices.md` - CI/CDパイプラインの使用方法

2. **技術スタック固有のガイド**
   - `references/framework-conventions.md` - 使用フレームワークの規約
   - `references/library-usage.md` - 主要ライブラリの使用パターン

3. **プロジェクト固有のパターン**
   - `references/error-handling.md` - エラーハンドリングパターン
   - `references/naming-conventions.md` - 命名規則の詳細

**ファイル作成の方針:**
- 各ファイルは単一の責務を持つようにする
- 具体的なコード例を含める
- プロジェクトの実装から抽出した実例を示す
- `project-convention.md` からは概要のみを記載し、詳細は個別ファイルで管理

## ステップ7: 中間ファイルの削除

分析に使用した中間ファイルを削除してクリーンアップします。

```bash
rm -f .claude/tmp/project-docs-analysis.md .claude/tmp/codebase-analysis.md
```

## ステップ8: 完了メッセージ

以下の情報を含む完了メッセージを表示します：

```
✓ {skill_name}スキルにプロジェクト固有の情報が適用されました

更新されたファイル:
  - plugins/{plugin_name}/skills/{skill_name}/SKILL.md
  - plugins/{plugin_name}/skills/{skill_name}/references/project-convention.md

作成されたreferenceファイル:
  [作成されたreferenceファイルのリスト]

これらのファイルには以下の情報が含まれています：
  - プロジェクトの開発規約
  - アーキテクチャパターン
  - コーディング標準
  - {skill_name}に関連するベストプラクティス
  - プロジェクト固有のツール使用方法

今後、{skill_name}スキルを使用する際は、これらのプロジェクト固有の規約に従ってください。
```

---

# 注意事項

## 前提条件
- スキル名は正確に指定してください（例：git-operations、figma-design-implementation）
- スキルが存在しない場合はエラーメッセージを表示
- プラグインが見つからない場合はエラーメッセージを表示

## 実行時の注意
- 分析には数分かかる場合があります
- 既存の `project-convention.md` は上書きされます（バックアップを推奨）
- 2つのエージェントは必ず並列で実行してください（効率化のため）
- **スキルの更新は必ずskill-creatorスキルを使用してください**
- referenceファイルは単一責務の原則に従って分割してください
- 既存のreferenceファイルがある場合は、内容を確認してから更新してください

## エラーハンドリング
- スキルが見つからない場合：「エラー: スキル '{skill_name}' が見つかりません」
- プラグインが特定できない場合：「エラー: スキルの所属プラグインを特定できません」
- 分析が失敗した場合：「エラー: 分析に失敗しました。詳細：[エラー内容]」

---

# 実行例

```bash
# git-operationsスキルにプロジェクト固有の情報を適用
/customize-skill git-operations

# figma-design-implementationスキルにプロジェクト固有の情報を適用
/customize-skill figma-design-implementation
```

---

それでは、{skill_name}スキルのカスタマイズを開始します。
