---
name: git-operation:branch
description: 現在の変更内容または命令内容に基づいてブランチを作成し、その後命令を遂行する。ユーザーが新しいブランチを作りたい場合、mainブランチで作業している場合、変更をブランチに切り出したい場合に使用する。
---

# Git Branch and Work

git-operations-expertサブエージェントでブランチを作成し、親エージェントが作業を実行する。

## サブエージェントへのプロンプト例

**作業内容が指定されている場合:**

```
以下の作業内容に基づいて適切なブランチ名を生成し、origin/mainから新しいブランチを作成してください。

作業内容: {ユーザーが指定した作業内容}

手順:
1. 作業内容を分析し、適切なtype(feature/fix/refactor/docs/testなど)を決定
2. ブランチ名を生成（形式: <type>/<brief-description>）
3. git fetch origin main を実行
4. git switch -c <ブランチ名> origin/main でブランチを作成
5. git branch --show-current で作成を確認

注意: ブランチ作成が完了したら終了。タスクの実装は行わないこと。
```

**変更内容がある場合:**

```
現在の変更内容を分析し、適切なブランチ名を生成してorigin/mainから新しいブランチを作成してください。

手順:
1. git status, git diff, git diff --staged を並列実行して変更内容を確認
2. 変更内容からtype(feature/fix/refactor/docs/testなど)を判断
3. ブランチ名を生成（形式: <type>/<brief-description>）
4. git fetch origin main を実行
5. git switch -c <ブランチ名> origin/main でブランチを作成
6. git branch --show-current で作成を確認

注意: ブランチ作成が完了したら終了。変更内容のステージングや追加作業は行わないこと。
```

**変更がなくヒアリングが必要な場合:**

```
ユーザーに作業内容をヒアリングし、それに基づいてブランチを作成してください。

手順:
1. git status で変更がないことを確認
2. AskUserQuestionツールでユーザーに作業内容をヒアリング
3. ヒアリング内容からtype(feature/fix/refactor/docs/testなど)を判断
4. ブランチ名を生成（形式: <type>/<brief-description>）
5. git fetch origin main を実行
6. git switch -c <ブランチ名> origin/main でブランチを作成
7. git branch --show-current で作成を確認

注意: ブランチ作成が完了したら終了。タスクの実装は行わないこと。
```

## ブランチ命名規則

```
<type>/<brief-description>
```

- **type**: feature, fix, refactor, docs, test など（fix/hotfixは明示的に指定されたときのみ使う）
- **brief-description**: 作業内容を簡潔に表す英語（kebab-case）

## 処理フロー

1. **サブエージェント起動**: Taskツールでgit-operations-expertを起動し、ブランチ作成のみ実行させる
2. **ブランチ作成完了**: サブエージェントからブランチ名と作成完了の報告を受ける
3. **作業実行（親エージェント）**: 以下を実行
   - 引数がある場合: その内容を実装
   - 変更が既にある場合: 変更を確認して必要に応じて追加作業
   - ヒアリングした場合: ヒアリング内容に基づいて作業を実行

## 注意事項

- 未コミットの変更は新しいブランチに引き継がれる
- 既に同名のブランチが存在する場合は異なる名前を提案する
