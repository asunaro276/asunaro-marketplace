# Git ベストプラクティス

このドキュメントは、Gitを効果的かつ安全に使用するためのベストプラクティスを提供する。

## コミット戦略

### コミットの粒度

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

### コミットメッセージ

#### 基本的な構造

```
<type>: <subject>

<body>

<footer>
```

#### Type（コミットの種類）

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

#### 例

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

#### ガイドライン

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

## ブランチ戦略

### Git Flow

```
main (production)
  └── develop (development)
       ├── feature/feature-name
       ├── bugfix/bug-name
       └── hotfix/hotfix-name
```

#### ブランチの種類

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

### GitHub Flow（シンプルな代替）

```
main (production)
  └── feature/feature-name
```

#### 特徴

1. **シンプル**: mainブランチのみ
2. **継続的デプロイ**: mainへのマージで自動デプロイ
3. **プルリクエスト**: すべての変更はPR経由
4. **レビュー**: マージ前にコードレビュー

### ブランチ命名規則

```
<type>/<description>-<issue-number>
```

#### 例

- `feature/user-authentication-123`
- `fix/login-error-456`
- `refactor/database-queries-789`
- `docs/api-documentation-101`
- `claude/add-git-skills-011CULfefegVmXpcz8RQKFVX`

## マージ戦略

### マージの種類

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

### マージ vs リベース

#### マージを使用する場合

- チームで共有しているブランチ
- 履歴を保持したい場合
- 安全性を優先する場合

#### リベースを使用する場合

- ローカルブランチ
- クリーンな履歴を作りたい場合
- まだプッシュしていないコミット

#### リベースの注意点

- **公開ブランチをリベースしない**
- **他の人が作業しているブランチをリベースしない**
- **履歴を書き換える操作であることを理解する**

## リモート操作

### プッシュのベストプラクティス

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

### プルのベストプラクティス

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

## セキュリティとプライバシー

### 認証情報の管理

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

### Git フック

#### pre-commit フック

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

#### commit-msg フック

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
git checkout <commit-hash>
git checkout -b <recovery-branch>
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

## パフォーマンス最適化

### リポジトリのクリーンアップ

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

### 浅いクローン

```bash
# 最新のコミットのみクローン
git clone --depth 1 <repository-url>

# 特定の深さでクローン
git clone --depth 10 <repository-url>

# 浅いクローンを完全な履歴に変換
git fetch --unshallow
```

### 部分クローン

```bash
# ブロブなしでクローン
git clone --filter=blob:none <repository-url>

# 大きなファイルを除外してクローン
git clone --filter=blob:limit=1m <repository-url>
```

## チーム協業

### プルリクエストのベストプラクティス

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

### コードレビューのガイドライン

#### レビュアーとして

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

#### 作成者として

1. **レビューを受け入れる姿勢**
   - フィードバックに感謝
   - 建設的に議論

2. **変更を反映**
   - コメントに対応
   - 理由を説明

3. **テストを追加**
   - 新機能にはテスト
   - バグ修正には再現テスト

## まとめ

- **小さく、頻繁にコミット**
- **明確で説明的なコミットメッセージ**
- **ブランチ戦略を一貫して適用**
- **マージ vs リベースを理解して使い分ける**
- **認証情報をコミットしない**
- **プッシュ前に確認**
- **ネットワークエラーはリトライ**
- **チームで協力してレビュー**
