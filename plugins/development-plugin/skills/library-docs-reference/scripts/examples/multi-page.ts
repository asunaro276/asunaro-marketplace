/**
 * 複数ページのドキュメントを取得する例
 *
 * Context7は1つのトピックについて最大10ページまでのドキュメントを
 * 提供します。このスクリプトは複数ページを効率的に取得する方法を示します。
 */

import { getMultiplePages, closeConnection } from '../context7-workflow.js';
import { writeFileSync } from 'fs';
import { resolve } from 'path';

async function main() {
  const libraryId = process.argv[2] || '/vercel/next.js';
  const topic = process.argv[3] || 'routing';
  const mode = (process.argv[4] || 'code') as 'code' | 'info';
  const maxPages = parseInt(process.argv[5] || '3', 10);

  console.log(`📚 複数ページのドキュメント取得例`);
  console.log(`   Library: ${libraryId}`);
  console.log(`   Topic: ${topic}`);
  console.log(`   Mode: ${mode}`);
  console.log(`   Pages: ${maxPages}\n`);

  try {
    // 複数ページを取得
    const allDocs = await getMultiplePages(libraryId, topic, mode, maxPages);

    console.log(`\n✅ ${allDocs.length} ページを取得しました\n`);

    // 各ページの内容を表示
    allDocs.forEach((doc, index) => {
      console.log(`${'='.repeat(80)}`);
      console.log(`ページ ${index + 1}/${allDocs.length}`);
      console.log(`${'='.repeat(80)}`);
      console.log(doc.content.substring(0, 1000));
      console.log('...\n');
    });

    // オプション: すべてのページを1つのファイルに保存
    const combinedContent = allDocs
      .map((doc, i) => `# Page ${i + 1}\n\n${doc.content}`)
      .join('\n\n' + '='.repeat(80) + '\n\n');

    const outputPath = resolve(
      process.cwd(),
      `docs-${libraryId.replace(/\//g, '-')}-${topic}-${mode}.md`
    );

    writeFileSync(outputPath, combinedContent, 'utf-8');
    console.log(`📄 ドキュメントを保存しました: ${outputPath}`);

  } catch (error) {
    console.error('❌ エラーが発生しました:', error);
    process.exit(1);
  } finally {
    await closeConnection();
  }
}

// スクリプトとして実行された場合
if (import.meta.url === `file://${process.argv[1]}`) {
  main();
}

/*
使用例:

# Next.jsのルーティングドキュメントを3ページ取得
tsx multi-page.ts /vercel/next.js routing code 3

# Reactのhooksドキュメントを5ページ取得
tsx multi-page.ts /facebook/react hooks code 5

# Vueの情報ドキュメントを2ページ取得
tsx multi-page.ts /vuejs/vue composition-api info 2
*/
