# タスク完了時のチェックリスト

## コード変更後の確認事項

### 1. ドキュメント更新
- [ ] SKILL.md ファイルが最新の内容を反映しているか
- [ ] README.md （存在する場合）が更新されているか
- [ ] コード内のコメントが正確か
- [ ] 使用例が動作することを確認したか

### 2. スキル開発（TypeScript/Node.js）

#### ビルドとテスト
```bash
# TypeScriptのビルド
cd scripts
npm run build

# 基本的な動作確認
npm run example:basic
```

#### コード品質
- [ ] 型定義が適切か
- [ ] エラーハンドリングが実装されているか
- [ ] リソースリークがないか（接続のクローズなど）
- [ ] 再利用可能な関数として設計されているか

### 3. MCP設定の確認

#### git-operation
- [ ] `.mcp.json` にContext7の設定が正しく記載されているか
- [ ] APIキーが有効か
- [ ] Serena MCPの設定が正しいか

#### frontend-plugin
- [ ] `plugin.json` にFigma Dev Mode MCPの設定があるか
- [ ] Serena MCPの設定が正しいか

### 4. プラグイン設定の確認

- [ ] `plugin.json` の必須フィールドがすべて記載されているか
  - name
  - description（日本語）
  - version
  - author
- [ ] スキルの `SKILL.md` にフロントマターが正しく記載されているか

### 5. Git操作

#### コミット前
```bash
# 変更内容の確認
git status
git diff

# 変更をステージング
git add <files>

# コミット
git commit -m "適切なコミットメッセージ"
```

#### コミットメッセージ
- [ ] 日本語または英語で明確に記述
- [ ] 変更の目的が分かるか
- [ ] 必要に応じて詳細説明を含めているか

#### プッシュ
```bash
# リモートにプッシュ
git push -u origin <branch-name>
```

### 6. プロジェクト固有の規約

#### language
- [ ] ユーザーとのやり取りは日本語で行っているか
- [ ] ドキュメントは日本語で記述されているか

#### ファイル配置
- [ ] プラグインは `plugins/` ディレクトリに配置されているか
- [ ] スキルは適切な `skills/` ディレクトリに配置されているか
- [ ] 設定ファイルは正しい場所に配置されているか

### 7. 動作確認

#### スキルの動作確認
```bash
# library-docs-reference の例
cd plugins/git-operation/skills/library-docs-reference/scripts
tsx context7-workflow.ts react hooks
```

#### エラーチェック
- [ ] 想定される入力でエラーが発生しないか
- [ ] エッジケースを考慮しているか
- [ ] エラーメッセージが分かりやすいか

### 8. 最終確認

- [ ] すべての変更が意図した通りに動作するか
- [ ] 既存の機能を壊していないか
- [ ] パフォーマンスに問題がないか
- [ ] セキュリティ上の問題がないか（APIキーの扱いなど）

## トラブルシューティング

### よくある問題

1. **MCP server not found**
   - `.mcp.json` または `plugin.json` の設定を確認
   - MCPサーバーが起動しているか確認

2. **型エラー (TypeScript)**
   ```bash
   npm install @types/node
   ```

3. **接続タイムアウト**
   - ネットワーク接続を確認
   - APIキーの有効性を確認

4. **依存関係のエラー**
   ```bash
   rm -rf node_modules package-lock.json
   npm install
   ```
