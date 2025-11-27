/**
 * 複数のライブラリを検索・比較する例
 *
 * このスクリプトは、複数のライブラリのドキュメントを並行して取得し、
 * 比較する方法を示します。
 */

import {
  resolveLibraryId,
  getDocsByLibraryName,
  closeConnection
} from '../context7-workflow.js';

interface LibraryInfo {
  name: string;
  id?: string;
  score?: number;
  snippetCount?: number;
  reputation?: string;
  docs?: string;
}

async function searchLibraries(libraryNames: string[]): Promise<LibraryInfo[]> {
  console.log(`🔍 ${libraryNames.length} 個のライブラリを検索中...\n`);

  const results = await Promise.all(
    libraryNames.map(async (name) => {
      try {
        const resolved = await resolveLibraryId(name);
        const selected = resolved.selected || resolved.matches?.[0];

        if (!selected) {
          return { name, id: undefined };
        }

        return {
          name,
          id: selected.id,
          score: selected.benchmark_score,
          snippetCount: selected.code_snippet_count,
          reputation: selected.source_reputation
        };
      } catch (error) {
        console.error(`❌ ${name} の検索に失敗:`, error);
        return { name, id: undefined };
      }
    })
  );

  return results;
}

async function compareLibraries(
  libraryNames: string[],
  topic: string,
  mode: 'code' | 'info' = 'code'
): Promise<void> {
  console.log(`📊 ライブラリ比較: ${libraryNames.join(' vs ')}`);
  console.log(`   トピック: ${topic}`);
  console.log(`   モード: ${mode}\n`);

  // まずライブラリ情報を取得
  const librariesInfo = await searchLibraries(libraryNames);

  // 比較表を表示
  console.log('=' .repeat(100));
  console.log('ライブラリ情報');
  console.log('='.repeat(100));
  console.log('Name'.padEnd(20), 'ID'.padEnd(30), 'Score'.padEnd(10), 'Snippets'.padEnd(10), 'Reputation');
  console.log('-'.repeat(100));

  librariesInfo.forEach(lib => {
    if (lib.id) {
      console.log(
        lib.name.padEnd(20),
        (lib.id || 'N/A').padEnd(30),
        (lib.score?.toString() || 'N/A').padEnd(10),
        (lib.snippetCount?.toString() || 'N/A').padEnd(10),
        lib.reputation || 'N/A'
      );
    } else {
      console.log(lib.name.padEnd(20), 'Not Found'.padEnd(30));
    }
  });
  console.log('='.repeat(100) + '\n');

  // 各ライブラリのドキュメントを取得
  console.log('📚 ドキュメントを取得中...\n');

  for (const lib of librariesInfo) {
    if (!lib.id) {
      console.log(`⚠️  ${lib.name}: ライブラリが見つかりませんでした\n`);
      continue;
    }

    try {
      const docs = await getDocsByLibraryName(lib.name, topic, mode, 1);

      console.log('='.repeat(100));
      console.log(`${lib.name} (${lib.id}) - ${topic}`);
      console.log('='.repeat(100));
      console.log(docs.content.substring(0, 800));
      console.log('...\n');

      lib.docs = docs.content;
    } catch (error) {
      console.error(`❌ ${lib.name} のドキュメント取得に失敗:`, error);
    }
  }

  console.log('✅ 比較完了！');
}

async function main() {
  const args = process.argv.slice(2);

  if (args.length < 2) {
    console.log(`
使用方法:
  tsx search-and-compare.ts <library1> <library2> [library3...] --topic <topic> [--mode <mode>]

例:
  # ReactとVueのhooksを比較
  tsx search-and-compare.ts react vue --topic hooks

  # Next.js、Nuxt、Remixのルーティングを比較
  tsx search-and-compare.ts next.js nuxt remix --topic routing --mode info

  # ExpressとFastifyのミドルウェアを比較
  tsx search-and-compare.ts express fastify --topic middleware --mode code
    `);
    process.exit(1);
  }

  // 引数をパース
  const topicIndex = args.indexOf('--topic');
  const modeIndex = args.indexOf('--mode');

  const libraries = args.slice(0, topicIndex !== -1 ? topicIndex : args.length);
  const topic = topicIndex !== -1 ? args[topicIndex + 1] : 'getting-started';
  const mode = (modeIndex !== -1 ? args[modeIndex + 1] : 'code') as 'code' | 'info';

  try {
    await compareLibraries(libraries, topic, mode);
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
