/**
 * Context7 ワークフローの基本的な使用例
 *
 * このスクリプトは、Context7を使用してライブラリのドキュメントを
 * 取得する基本的な方法を示します。
 */

import {
  resolveLibraryId,
  getLibraryDocs,
  getDocsByLibraryName,
  closeConnection
} from '../context7-workflow.js';

async function main() {
  try {
    console.log('🚀 Context7 基本的な使用例\n');

    // 例1: ライブラリIDを解決
    console.log('例1: ライブラリIDの解決');
    console.log('-'.repeat(80));
    const reactResolve = await resolveLibraryId('react');
    console.log('マッチ数:', reactResolve.matches?.length);
    console.log('選択されたライブラリ:', reactResolve.selected);
    console.log('\n');

    // 例2: ライブラリIDから直接ドキュメントを取得
    console.log('例2: ライブラリIDからドキュメントを取得');
    console.log('-'.repeat(80));
    if (reactResolve.selected) {
      const docs = await getLibraryDocs(
        reactResolve.selected.id,
        'hooks',
        'code',
        1
      );
      console.log('取得したドキュメント:');
      console.log(docs.content.substring(0, 500) + '...\n');
    }

    // 例3: ライブラリ名から直接ドキュメントを取得（推奨方法）
    console.log('例3: ライブラリ名から直接ドキュメントを取得');
    console.log('-'.repeat(80));
    const vueDocs = await getDocsByLibraryName('vue', 'reactive', 'code', 1);
    console.log('取得したドキュメント:');
    console.log(vueDocs.content.substring(0, 500) + '...\n');

    // 例4: 情報モードでドキュメントを取得
    console.log('例4: 情報モードでドキュメントを取得');
    console.log('-'.repeat(80));
    const nextDocs = await getDocsByLibraryName('next.js', 'routing', 'info', 1);
    console.log('取得したドキュメント:');
    console.log(nextDocs.content.substring(0, 500) + '...\n');

    console.log('✅ すべての例が正常に完了しました！');

  } catch (error) {
    console.error('❌ エラーが発生しました:', error);
    process.exit(1);
  } finally {
    // 接続をクローズ
    await closeConnection();
  }
}

// スクリプトとして実行された場合
if (import.meta.url === `file://${process.argv[1]}`) {
  main();
}
