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
- ブランチ名は `feature/` で始まり、セッションIDで終わること（そうでない場合403エラー）
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

### コミット戦略

#### コミットの粒度

1. **小さく、論理的な単位でコミット**
   - 1つのコミットは1つの論理的な変更を表す
   - 複数の無関係な変更を1つのコミットにまとめない
   - 大きな機能は複数の小さなコミットに分割

2. **頻繁にコミット**
   - 作業を小さな単位で保存
   - 問題が発生した場合に簡単にロールバック可能
   - チームメンバーとの衝突を最小化

3. **完全な状態でコミット**
   - ビルドが通る状態でコミット
   - テストが通る状態でコミット
   - 部分的な実装はコミットしない（または WIP とマーク）

#### コミットメッセージ

##### 基本的な構造

```
<type>: <subject>

<body>

<footer>
```

##### Type（コミットの種類）

- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響しない変更（フォーマット、セミコロンなど）
- `refactor`: バグ修正や機能追加ではないコードの変更
- `perf`: パフォーマンス改善
- `test`: テストの追加や修正
- `build`: ビルドシステムや外部依存関係の変更
- `ci`: CI設定ファイルやスクリプトの変更
- `chore`: その他の変更（ソースコードやテストに影響しない変更）

##### 例

```
feat: ユーザー認証機能を追加

JWT を使用したユーザー認証システムを実装。
ログイン、ログアウト、トークンのリフレッシュ機能を含む。

Closes #123
```

```
fix: ログインフォームのバリデーションエラーを修正

メールアドレスの形式チェックが正しく動作していなかった問題を修正。
正規表現を更新し、より厳密な検証を実装。

Fixes #456
```

##### ガイドライン

1. **subject（件名）**
   - 50文字以内
   - 命令形を使用（「追加した」ではなく「追加」）
   - 最初の文字は大文字
   - 末尾にピリオドを付けない

2. **body（本文）**
   - 72文字で折り返し
   - 「何を」ではなく「なぜ」を説明
   - 変更の理由と背景を記述
   - subjectとbodyの間に空行を入れる

3. **footer（フッター）**
   - 関連するイシュー番号を記載
   - Breaking changes を記載
   - Co-authored-by を記載（ペアプログラミングの場合）

### ブランチ戦略

#### Git Flow

```
main (production)
  └── develop (development)
       ├── feature/feature-name
       ├── bugfix/bug-name
       └── hotfix/hotfix-name
```

##### ブランチの種類

1. **main/master**: 本番環境のコード
   - 常に安定した状態
   - 直接コミットしない
   - タグ付けしてバージョン管理

2. **develop**: 開発ブランチ
   - 次のリリースの統合ブランチ
   - 機能ブランチをマージ

3. **feature/**: 機能ブランチ
   - developから分岐
   - developにマージ
   - 命名: `feature/user-authentication`

4. **bugfix/**: バグ修正ブランチ
   - developから分岐
   - developにマージ
   - 命名: `bugfix/login-error`

5. **hotfix/**: 緊急修正ブランチ
   - mainから分岐
   - mainとdevelopにマージ
   - 命名: `hotfix/security-patch`

6. **release/**: リリースブランチ
   - developから分岐
   - mainとdevelopにマージ
   - 命名: `release/v1.2.0`

#### GitHub Flow（シンプルな代替）

```
main (production)
  └── feature/feature-name
```

##### 特徴

1. **シンプル**: mainブランチのみ
2. **継続的デプロイ**: mainへのマージで自動デプロイ
3. **プルリクエスト**: すべての変更はPR経由
4. **レビュー**: マージ前にコードレビュー

#### ブランチ命名規則

```
<type>/<description>-<issue-number>
```

##### 例

- `feature/user-authentication-123`
- `fix/login-error-456`
- `refactor/database-queries-789`
- `docs/api-documentation-101`
- `claude/add-git-skills-011CULfefegVmXpcz8RQKFVX`

### マージ戦略

#### マージの種類

1. **Fast-forward マージ**
   ```bash
   git merge <branch-name>
   ```
   - 履歴が一直線
   - マージコミットなし
   - 適用: シンプルな変更

2. **No-fast-forward マージ**
   ```bash
   git merge --no-ff <branch-name>
   ```
   - 常にマージコミットを作成
   - 機能の履歴を保持
   - 適用: 機能ブランチのマージ

3. **Squash マージ**
   ```bash
   git merge --squash <branch-name>
   ```
   - すべてのコミットを1つにまとめる
   - クリーンな履歴
   - 適用: 細かいコミットが多い場合

4. **Rebase マージ**
   ```bash
   git rebase <base-branch>
   ```
   - 履歴を一直線に
   - コミットを再適用
   - 適用: 履歴をクリーンに保ちたい場合

#### マージ vs リベース

##### マージを使用する場合

- チームで共有しているブランチ
- 履歴を保持したい場合
- 安全性を優先する場合

##### リベースを使用する場合

- ローカルブランチ
- クリーンな履歴を作りたい場合
- まだプッシュしていないコミット

##### リベースの注意点

- **公開ブランチをリベースしない**
- **他の人が作業しているブランチをリベースしない**
- **履歴を書き換える操作であることを理解する**

### リモート操作のベストプラクティス

#### プッシュのベストプラクティス

1. **プッシュ前に確認**
   ```bash
   git log origin/<branch-name>..HEAD
   git diff origin/<branch-name>..HEAD
   ```

2. **上流ブランチを設定**
   ```bash
   git push -u origin <branch-name>
   ```

3. **強制プッシュは慎重に**
   ```bash
   # 推奨: より安全
   git push --force-with-lease

   # 非推奨: 危険
   git push --force
   ```

4. **ネットワークエラー時はリトライ**
   ```bash
   for i in 1 2 3 4; do
     git push -u origin <branch-name> && break
     sleep $((2 ** i))
   done
   ```

#### プルのベストプラクティス

1. **作業前に最新化**
   ```bash
   git pull origin <branch-name>
   ```

2. **リベースしながらプル**
   ```bash
   git pull --rebase origin <branch-name>
   ```

3. **特定のブランチのみフェッチ**
   ```bash
   git fetch origin <branch-name>
   ```

4. **ネットワークエラー時はリトライ**
   ```bash
   for i in 1 2 3 4; do
     git fetch origin <branch-name> && break
     sleep $((2 ** i))
   done
   ```

### セキュリティとプライバシー

#### 認証情報の管理

1. **認証情報をコミットしない**
   - `.env` ファイル
   - API キー
   - パスワード
   - 秘密鍵
   - トークン

2. **`.gitignore` を活用**
   ```gitignore
   # 環境変数
   .env
   .env.local
   .env.*.local

   # 認証情報
   credentials.json
   secrets.yaml
   *.pem
   *.key

   # IDE設定（個人設定を含む）
   .vscode/
   .idea/
   ```

3. **誤ってコミットした場合**
   ```bash
   # 履歴から削除（すべてのブランチから）
   git filter-branch --force --index-filter \
     'git rm --cached --ignore-unmatch <file>' \
     --prune-empty --tag-name-filter cat -- --all

   # または git-filter-repo を使用（推奨）
   git filter-repo --path <file> --invert-paths
   ```

#### Git フック

##### pre-commit フック

```bash
#!/bin/sh
# .git/hooks/pre-commit

# 認証情報のチェック
if git diff --cached | grep -E 'password|secret|api_key'; then
  echo "Error: Potential credentials found!"
  exit 1
fi

# コードフォーマットのチェック
make fmt-check
```

##### commit-msg フック

```bash
#!/bin/sh
# .git/hooks/commit-msg

# コミットメッセージの形式チェック
commit_msg=$(cat "$1")
if ! echo "$commit_msg" | grep -qE '^(feat|fix|docs|style|refactor|perf|test|build|ci|chore):'; then
  echo "Error: Commit message must start with a type!"
  exit 1
fi
```

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

### パフォーマンス最適化

#### リポジトリのクリーンアップ

```bash
# ガベージコレクション
git gc

# 積極的なガベージコレクション
git gc --aggressive --prune=now

# リモートで削除されたブランチを削除
git remote prune origin

# 未追跡ファイルを削除
git clean -fd
```

#### 浅いクローン

```bash
# 最新のコミットのみクローン
git clone --depth 1 <repository-url>

# 特定の深さでクローン
git clone --depth 10 <repository-url>

# 浅いクローンを完全な履歴に変換
git fetch --unshallow
```

#### 部分クローン

```bash
# ブロブなしでクローン
git clone --filter=blob:none <repository-url>

# 大きなファイルを除外してクローン
git clone --filter=blob:limit=1m <repository-url>
```

### チーム協業

#### プルリクエストのベストプラクティス

1. **小さく保つ**
   - 200〜400行が理想
   - 1つの機能や修正に集中

2. **説明的なタイトルと説明**
   - 何を変更したか
   - なぜ変更したか
   - どのようにテストしたか

3. **レビュアーを指定**
   - 関連する専門知識を持つ人
   - 影響を受けるコンポーネントの担当者

4. **CI/CD を通過させる**
   - すべてのテストが通る
   - ビルドが成功する
   - リンターが通る

#### コードレビューのガイドライン

##### レビュアーとして

1. **建設的なフィードバック**
   - 問題点を指摘するだけでなく、解決策を提案
   - 肯定的なコメントも忘れずに

2. **優先順位を明確に**
   - 必須: 修正が必要
   - 推奨: 改善の余地がある
   - 提案: 検討してほしい

3. **迅速に対応**
   - 24時間以内にレビュー
   - ブロッキングな問題は優先

##### 作成者として

1. **レビューを受け入れる姿勢**
   - フィードバックに感謝
   - 建設的に議論

2. **変更を反映**
   - コメントに対応
   - 理由を説明

3. **テストを追加**
   - 新機能にはテスト
   - バグ修正には再現テスト

## エラーハンドリング

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

## トラブルシューティング

### よくある問題と解決方法

#### 1. マージコンフリクト

```bash
# コンフリクトを確認
git status

# コンフリクトを手動で解決
# エディタでファイルを編集

# 解決したファイルをステージング
git add <resolved-file>

# マージを完了
git commit
```

#### 2. 誤ったコミット

```bash
# 直前のコミットを修正
git commit --amend

# コミットを取り消し（変更は保持）
git reset --soft HEAD~1

# コミットを取り消し（変更も破棄）
git reset --hard HEAD~1
```

#### 3. 削除されたコミットの復元

```bash
# reflog で削除されたコミットを検索
git reflog

# コミットを復元
git restore <commit-hash>
git switch -c <recovery-branch>
```

#### 4. リベースの中断

```bash
# リベースを中止
git rebase --abort

# リベースをスキップ
git rebase --skip

# リベースを続行
git rebase --continue
```

## まとめ

- **小さく、頻繁にコミット**
- **明確で説明的なコミットメッセージ**
- **ブランチ戦略を一貫して適用**
- **マージ vs リベースを理解して使い分ける**
- **認証情報をコミットしない**
- **プッシュ前に確認**
- **ネットワークエラーはリトライ**
- **チームで協力してレビュー**
