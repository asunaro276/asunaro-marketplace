---
name: git-operation:create-pr
description: 現在のブランチの変更内容を分析し、PRテンプレートに従った形式でドラフトプルリクエストを作成する。ユーザーがPRを作成したい場合、プルリクエストを出したい場合に使用する。
---

# Git Create PR

git-operations-expertサブエージェントを使用して、現在ブランチの変更を分析しPRテンプレートに従ったドラフトPRを作成する。

## サブエージェントへのプロンプト例

```
現在のブランチの変更内容を分析して、ドラフトプルリクエストを作成してください。

以下の手順で実行してください：

1. 現在の状態確認（git status, git diff --staged, git diff, git status -sb, git log main...HEAD --oneline, git diff main...HEAD を並列実行）

2. 前提条件の確認:
  - mainブランチでないことを確認（main/masterの場合は `/git-branch` の実行を促して終了）
  - コミットが存在することを確認（コミットがない場合は `/git-commit` の実行を促して終了）

3. PRテンプレートの確認:
  - `git rev-parse --show-toplevel` でGitルートディレクトリを取得する
  - Gitルートディレクトリ配下の以下のパスを順番に探索する:
    - .github/PULL_REQUEST_TEMPLATE.md
    - .github/pull_request_template.md
    - .github/PULL_REQUEST_TEMPLATE/ 配下の最初のファイル
    - PULL_REQUEST_TEMPLATE.md

4. すべてのコミット履歴と変更内容を分析してPR内容を作成:
  - タイトル: 変更内容を総括した簡潔な日本語タイトル（1行）
  - 本文: テンプレートがあればそれに従い、なければ以下の形式:
    ```
    ## Summary
    <変更内容を1〜3つの箇条書きで要約>

    ## Changes
    <主要な変更点を箇条書きで記載>

    ## Test plan
    - [ ] テスト項目1

    🤖 Generated with [Claude Code](https://claude.com/claude-code)
    ```

5. ベースブランチの決定（gh repo view --json defaultBranchRef --jq .defaultBranchRef.name）

6. 必要に応じてリモートにプッシュ:
  - git push -u origin <ブランチ名>
  - ネットワークエラー時は指数バックオフでリトライ（最大4回: 2s, 4s, 8s, 16s）

7. ドラフトPRを作成:
  gh pr create --draft --title "<タイトル>" --body "$(cat <<'EOF'
  <本文>
  EOF
  )"

8. 作成されたPRのURLを返す

9. /loop を使って5分毎にCIを監視する:
  `/loop 5m gh pr checks <PR番号> --watch`
  - CIが成功/失敗で完了したら監視を終了する
```

## 注意事項

- main/masterブランチからは実行不可
- コミットがない場合は `/git-commit` を先に実行する必要がある
- リモートにプッシュされていない変更は自動的にプッシュされる
- PRはドラフト状態で作成される（レビュー準備ができたら手動でReady for reviewに変更）
- PRが既に存在する場合は既存PRのURLを表示して終了
