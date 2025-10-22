# Git コマンドリファレンス

このドキュメントは、頻繁に使用されるGitコマンドの包括的なリファレンスを提供する。

## 基本コマンド

### 初期化とクローン

```bash
# 新しいGitリポジトリを初期化
git init

# リポジトリをクローン
git clone <repository-url>

# 特定のブランチをクローン
git clone -b <branch-name> <repository-url>

# 浅いクローン（履歴を制限）
git clone --depth 1 <repository-url>
```

### 設定

```bash
# ユーザー名を設定
git config --global user.name "Your Name"

# メールアドレスを設定
git config --global user.email "your.email@example.com"

# 設定を表示
git config --list

# 特定の設定を表示
git config user.name
```

## ステージングとコミット

### 変更の追加

```bash
# 特定のファイルをステージング
git add <file>

# 特定のパターンのファイルをステージング
git add *.js

# インタラクティブにステージング（自動化環境では非推奨）
# git add -i

# パッチモードでステージング（自動化環境では非推奨）
# git add -p
```

### コミット

```bash
# コミットメッセージ付きでコミット
git commit -m "commit message"

# 詳細なコミットメッセージ（エディタが開く）
git commit

# ステージングとコミットを同時に（追跡ファイルのみ）
git commit -am "commit message"

# 直前のコミットを修正
git commit --amend -m "new message"

# ファイルを追加して直前のコミットを修正
git add <file>
git commit --amend --no-edit

# 空のコミットを作成
git commit --allow-empty -m "empty commit"
```

### 変更の取り消し

```bash
# ファイルのステージングを解除
git restore --staged <file>

# ファイルの変更を破棄
git restore <file>

# すべての変更を破棄
git restore .

# 直前のコミットを取り消し（変更は保持）
git reset --soft HEAD~1

# 直前のコミットを取り消し（ステージングも解除）
git reset --mixed HEAD~1

# 直前のコミットを取り消し（変更も破棄）
git reset --hard HEAD~1

# 特定のコミットに戻る
git reset --hard <commit-hash>

# コミットを打ち消す新しいコミットを作成
git revert <commit-hash>
```

## ブランチ操作

### ブランチの作成と切り替え

```bash
# ブランチを作成
git branch <branch-name>

# ブランチに切り替え
git checkout <branch-name>

# ブランチを作成して切り替え
git checkout -b <branch-name>

# 特定のコミットからブランチを作成
git checkout -b <branch-name> <commit-hash>

# リモートブランチから新しいブランチを作成
git checkout -b <branch-name> origin/<remote-branch>

# 前のブランチに戻る
git checkout -
```

### ブランチの管理

```bash
# ローカルブランチを表示
git branch

# すべてのブランチを表示
git branch -a

# リモートブランチを表示
git branch -r

# 現在のブランチを表示
git branch --show-current

# ブランチ名を変更
git branch -m <old-name> <new-name>

# 現在のブランチ名を変更
git branch -m <new-name>

# ブランチを削除
git branch -d <branch-name>

# 強制削除（マージされていない変更がある場合）
git branch -D <branch-name>

# マージ済みのブランチを表示
git branch --merged

# 未マージのブランチを表示
git branch --no-merged
```

## リモート操作

### リモートの管理

```bash
# リモートを表示
git remote

# リモートの詳細を表示
git remote -v

# リモートを追加
git remote add <name> <url>

# リモートを削除
git remote remove <name>

# リモート名を変更
git remote rename <old-name> <new-name>

# リモートのURLを変更
git remote set-url <name> <new-url>

# リモートの情報を表示
git remote show <name>
```

### フェッチとプル

```bash
# リモートから変更をフェッチ
git fetch

# 特定のリモートからフェッチ
git fetch <remote-name>

# 特定のブランチをフェッチ
git fetch origin <branch-name>

# すべてのリモートからフェッチ
git fetch --all

# 削除されたリモートブランチを削除
git fetch --prune

# リモートから変更をプル（フェッチ＋マージ）
git pull

# 特定のブランチをプル
git pull origin <branch-name>

# リベースしながらプル
git pull --rebase

# すべてのリモートからプル
git pull --all
```

### プッシュ

```bash
# リモートにプッシュ
git push

# 特定のブランチをプッシュ
git push origin <branch-name>

# 上流ブランチを設定してプッシュ
git push -u origin <branch-name>

# すべてのブランチをプッシュ
git push --all

# タグをプッシュ
git push --tags

# リモートブランチを削除
git push origin --delete <branch-name>

# 強制プッシュ（注意: 破壊的操作）
git push --force

# より安全な強制プッシュ
git push --force-with-lease
```

## マージとリベース

### マージ

```bash
# ブランチをマージ
git merge <branch-name>

# Fast-forwardなしでマージ
git merge --no-ff <branch-name>

# マージコミットメッセージを指定
git merge -m "merge message" <branch-name>

# マージの中止
git merge --abort

# マージの続行（コンフリクト解決後）
git merge --continue

# スカッシュマージ
git merge --squash <branch-name>
```

### リベース

```bash
# ブランチをリベース
git rebase <base-branch>

# インタラクティブリベース（自動化環境では非推奨）
# git rebase -i <base-branch>

# リベースの続行
git rebase --continue

# リベースのスキップ
git rebase --skip

# リベースの中止
git rebase --abort

# 特定のコミットからリベース
git rebase --onto <new-base> <old-base> <branch>
```

## 履歴の確認

### ログ

```bash
# コミット履歴を表示
git log

# 簡潔な履歴表示
git log --oneline

# グラフ表示
git log --graph --oneline --all

# 特定の数のコミットを表示
git log -n 5

# 特定の期間のコミットを表示
git log --since="2 weeks ago"
git log --until="2023-01-01"

# 特定の作成者のコミットを表示
git log --author="Author Name"

# 特定のファイルの履歴を表示
git log -- <file>

# コミットメッセージを検索
git log --grep="search term"

# 変更内容を含めて表示
git log -p

# 統計情報を表示
git log --stat
```

### 差分

```bash
# 作業ツリーの変更を表示
git diff

# ステージングされた変更を表示
git diff --staged
git diff --cached

# 2つのコミット間の差分
git diff <commit1> <commit2>

# 2つのブランチ間の差分
git diff <branch1> <branch2>

# 特定のファイルの差分
git diff <file>

# 変更されたファイル名のみ表示
git diff --name-only

# 統計情報のみ表示
git diff --stat
```

### 検査

```bash
# コミットの詳細を表示
git show <commit-hash>

# ファイルの各行の最終変更を表示
git blame <file>

# ファイルの内容を特定のコミットで表示
git show <commit-hash>:<file>

# リファレンスのログを表示
git reflog
```

## スタッシュ

```bash
# 変更をスタッシュに保存
git stash

# メッセージ付きでスタッシュ
git stash save "message"

# 未追跡ファイルも含めてスタッシュ
git stash -u

# スタッシュリストを表示
git stash list

# スタッシュを適用（保持）
git stash apply

# 最新のスタッシュを適用して削除
git stash pop

# 特定のスタッシュを適用
git stash apply stash@{0}

# スタッシュの変更内容を表示
git stash show

# スタッシュの詳細な変更内容を表示
git stash show -p

# スタッシュを削除
git stash drop

# 特定のスタッシュを削除
git stash drop stash@{0}

# すべてのスタッシュを削除
git stash clear
```

## タグ

```bash
# タグを作成
git tag <tag-name>

# 注釈付きタグを作成
git tag -a <tag-name> -m "tag message"

# 特定のコミットにタグを付ける
git tag <tag-name> <commit-hash>

# タグを表示
git tag

# タグの詳細を表示
git show <tag-name>

# タグを削除
git tag -d <tag-name>

# リモートのタグを削除
git push origin --delete <tag-name>

# タグをプッシュ
git push origin <tag-name>

# すべてのタグをプッシュ
git push --tags
```

## クリーンアップ

```bash
# 未追跡ファイルを表示
git clean -n

# 未追跡ファイルを削除
git clean -f

# 未追跡ファイルとディレクトリを削除
git clean -fd

# 無視されたファイルも削除
git clean -fdx

# ガベージコレクション
git gc

# 積極的なガベージコレクション
git gc --aggressive

# リモートで削除されたブランチを削除
git remote prune origin
```

## サブモジュール

```bash
# サブモジュールを追加
git submodule add <repository-url> <path>

# サブモジュールを初期化
git submodule init

# サブモジュールを更新
git submodule update

# サブモジュールを初期化して更新
git submodule update --init

# すべてのサブモジュールを再帰的に更新
git submodule update --init --recursive

# サブモジュールの状態を表示
git submodule status
```

## 高度な操作

### Cherry-pick

```bash
# 特定のコミットを適用
git cherry-pick <commit-hash>

# 複数のコミットを適用
git cherry-pick <commit1> <commit2>

# Cherry-pickの中止
git cherry-pick --abort

# Cherry-pickの続行
git cherry-pick --continue
```

### Bisect

```bash
# Bisectを開始
git bisect start

# 悪いコミットをマーク
git bisect bad

# 良いコミットをマーク
git bisect good <commit-hash>

# Bisectをリセット
git bisect reset
```

### Archive

```bash
# リポジトリをアーカイブ
git archive --format=zip --output=archive.zip HEAD

# 特定のブランチをアーカイブ
git archive --format=tar --output=archive.tar <branch-name>
```

## トラブルシューティング

### 診断

```bash
# リポジトリの整合性をチェック
git fsck

# 統計情報を表示
git count-objects -v

# すべての設定を表示
git config --list --show-origin
```

### リカバリ

```bash
# 削除されたコミットを復元
git reflog
git checkout <commit-hash>

# 削除されたブランチを復元
git reflog
git checkout -b <branch-name> <commit-hash>

# 失われたコミットを見つける
git fsck --lost-found
```
