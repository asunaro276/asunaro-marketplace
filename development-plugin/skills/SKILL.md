---
name: git-operations
description: Git操作に関する包括的なガイド。ユーザーがブランチ作成、コミット、プッシュ、プル、マージ、リベースなどのGit操作を実行したい場合や、Gitのベストプラクティスに従う必要がある場合に使用すべきスキルです。
---

# Git Operations

## 概要

このスキルは、Git操作の実行に関する包括的なガイダンスを提供する。基本的なコミットとプッシュから、高度なブランチ管理とリモート操作まで、すべてのGit操作をカバーする。

## クイックスタート

### 基本的なワークフロー

1. **変更をステージング**: `git add <file>`
2. **変更をコミット**: `git commit -m "commit message"`
3. **リモートにプッシュ**: `git push -u origin <branch-name>`

## コア操作

### 1. ブランチ操作

#### ブランチの作成と切り替え

```bash
# 新しいブランチを作成
git branch <branch-name>

# ブランチに切り替え
git switch <branch-name>

# ブランチを作成して切り替え（ワンステップ）
git switch -c <branch-name>
```

#### ブランチの確認

```bash
# ローカルブランチを表示
git branch

# すべてのブランチ（リモート含む）を表示
git branch -a

# 現在のブランチを確認
git branch --show-current
```

#### ブランチの削除

```bash
# ローカルブランチを削除
git branch -d <branch-name>

# 強制削除（マージされていない変更がある場合）
git branch -D <branch-name>

# リモートブランチを削除
git push origin --delete <branch-name>
```

### 2. コミット操作

#### 基本的なコミット

```bash
# 変更をステージング
git add <file>

# コミットメッセージ付きでコミット
git commit -m "commit message"

# 詳細なコミットメッセージ（エディタが開く）
git commit
```

#### コミットの修正

```bash
# 直前のコミットを修正（メッセージを変更）
git commit --amend -m "new message"

# 直前のコミットにファイルを追加
git add <file>
git commit --amend --no-edit
```

**注意**: `--amend` は他の開発者のコミットには使用しないこと。必ず以下を確認：
- 作成者が自分であること: `git log -1 --format='%an %ae'`
- プッシュされていないこと: `git status` で「Your branch is ahead」を確認

#### コミット履歴の確認

```bash
# コミット履歴を表示
git log

# 簡潔な履歴表示
git log --oneline

# 特定の数のコミットを表示
git log -n 5

# グラフ表示
git log --graph --oneline --all
```

### 3. リモート操作

#### リモートの管理

```bash
# リモートを表示
git remote -v

# リモートを追加
git remote add <name> <url>

# リモートを削除
git remote remove <name>
```

#### フェッチとプル

```bash
# 特定のブランチをフェッチ（推奨）
git fetch origin <branch-name>

# すべてのリモートブランチをフェッチ
git fetch --all

# プル（フェッチ＋マージ）
git pull origin <branch-name>
```

**ネットワークエラー時のリトライロジック**:
```bash
# フェッチのリトライ（最大4回、指数バックオフ: 2s, 4s, 8s, 16s）
for i in 1 2 3 4; do
  git fetch origin <branch-name> && break
  sleep $((2 ** i))
done
```

#### プッシュ

```bash
# 基本的なプッシュ
git push origin <branch-name>

# 上流ブランチを設定してプッシュ（初回）
git push -u origin <branch-name>

# すべてのブランチをプッシュ（非推奨）
git push --all
```

**重要なプッシュのルール**:
- ブランチ名は `claude/` で始まり、セッションIDで終わること（そうでない場合403エラー）
- ネットワークエラー時は最大4回リトライ（指数バックオフ: 2s, 4s, 8s, 16s）
- `--force` は main/master ブランチには使用しない

**ネットワークエラー時のリトライロジック**:
```bash
# プッシュのリトライ（最大4回、指数バックオフ）
for i in 1 2 3 4; do
  git push -u origin <branch-name> && break
  sleep $((2 ** i))
done
```

### 4. 状態の確認

#### 作業ツリーの状態

```bash
# 状態を確認
git status

# 簡潔な状態表示
git status -s

# 変更内容を確認
git diff

# ステージングされた変更を確認
git diff --staged
```

#### ブランチ間の差分

```bash
# 2つのブランチの差分
git diff <branch1>..<branch2>

# 現在のブランチとメインブランチの差分
git diff main...HEAD

# 特定のブランチから分岐した後のコミット履歴
git log <base-branch>...HEAD
```

### 5. マージとリベース

#### マージ

```bash
# ブランチをマージ
git merge <branch-name>

# Fast-forwardなしでマージ
git merge --no-ff <branch-name>

# マージの中止
git merge --abort
```

#### リベース

```bash
# ブランチをリベース
git rebase <base-branch>

# リベースの続行
git rebase --continue

# リベースの中止
git rebase --abort
```

**注意**: インタラクティブモード（`-i` フラグ）は使用しないこと。`git rebase -i` や `git add -i` などは対話的な入力が必要なため、自動化環境ではサポートされていない。

### 6. スタッシュ

```bash
# 変更をスタッシュに保存
git stash

# スタッシュにメッセージを付けて保存
git stash save "message"

# スタッシュリストを表示
git stash list

# スタッシュを適用（保持）
git stash apply

# スタッシュを適用して削除
git stash pop

# 特定のスタッシュを適用
git stash apply stash@{0}

# スタッシュを削除
git stash drop
```

## ベストプラクティス

### コミットメッセージ

1. **明確で説明的**: コミットメッセージは「何を」ではなく「なぜ」を重視
2. **簡潔に**: 1〜2文で要約
3. **変更の性質を明確に**: 新機能、バグ修正、リファクタリングなどを明記
4. **リポジトリのスタイルに従う**: `git log` で既存のコミットメッセージスタイルを確認

### ブランチ命名規則

- **feature/**: 新機能用（例: `feature/user-authentication`）
- **fix/**: バグ修正用（例: `fix/login-error`）
- **hotfix/**: 緊急修正用
- **refactor/**: リファクタリング用
- **docs/**: ドキュメント更新用
- **claude/**: Claude Code用（例: `claude/add-feature-{session-id}`）

### Git安全性プロトコル

1. **設定を更新しない**: `git config` は変更しない
2. **破壊的コマンドを避ける**: `--force`, `--hard reset` などは明示的な指示がない限り使用しない
3. **フックをスキップしない**: `--no-verify`, `--no-gpg-sign` などは明示的な指示がない限り使用しない
4. **main/master への強制プッシュを避ける**: ユーザーが明示的に要求した場合のみ実行（警告を表示）
5. **コミット前に確認**: `git status` と `git diff` で変更内容を確認
6. **認証情報をコミットしない**: `.env`, `credentials.json` などのファイルは除外

### ワークフロー

1. **作業前に最新状態に**: `git pull origin <branch-name>`
2. **頻繁にコミット**: 小さな論理的な単位でコミット（`git add .`は使用不可）
3. **プッシュ前に確認**: `git log` で変更内容を確認
4. **ブランチで作業**: mainブランチに直接コミットしない
5. **定期的にプッシュ**: 作業を失わないように定期的にリモートにプッシュ

## エラーハンドリング

### ネットワークエラー

ネットワークエラーが発生した場合、指数バックオフでリトライする：
- 1回目: 2秒待機
- 2回目: 4秒待機
- 3回目: 8秒待機
- 4回目: 16秒待機

### マージコンフリクト

```bash
# コンフリクトを確認
git status

# コンフリクトを手動で解決後
git add <resolved-file>
git commit

# マージを中止
git merge --abort
```

### 誤ったコミットの修正

```bash
# 直前のコミットを取り消し（変更は保持）
git reset --soft HEAD~1

# 直前のコミットを取り消し（変更も破棄）
git reset --hard HEAD~1  # 注意: 破壊的操作

# 特定のコミットに戻る
git revert <commit-hash>
```

## リファレンス

詳細なコマンドリファレンスとベストプラクティスについては、`references/` ディレクトリを参照：
- `git_commands.md`: 完全なGitコマンドリファレンス
- `best_practices.md`: Gitのベストプラクティスと推奨事項
