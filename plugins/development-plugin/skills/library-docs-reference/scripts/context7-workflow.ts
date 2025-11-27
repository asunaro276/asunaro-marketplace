import { callMCPTool, closeConnection } from './mcp-client.js';

interface LibraryMatch {
  id: string;
  name: string;
  description?: string;
  benchmark_score?: number;
  code_snippet_count?: number;
  source_reputation?: string;
}

interface ResolveLibraryResponse {
  matches: LibraryMatch[];
  selected?: LibraryMatch;
}

interface Documentation {
  content: string;
  page?: number;
  hasMore?: boolean;
}

/**
 * ライブラリ名からContext7互換のライブラリIDを解決
 * @param libraryName 検索するライブラリ名
 * @returns 解決されたライブラリIDと候補リスト
 */
export async function resolveLibraryId(libraryName: string): Promise<ResolveLibraryResponse> {
  console.log(`\n📚 Resolving library ID for: ${libraryName}`);

  const result = await callMCPTool<ResolveLibraryResponse>(
    'resolve-library-id',
    { libraryName }
  );

  console.log(`✅ Found ${result.matches?.length || 0} matches`);
  if (result.selected) {
    console.log(`   Selected: ${result.selected.id} (${result.selected.name})`);
    if (result.selected.benchmark_score) {
      console.log(`   Quality Score: ${result.selected.benchmark_score}/100`);
    }
  }

  return result;
}

/**
 * ライブラリIDからドキュメントを取得
 * @param libraryId Context7互換のライブラリID (例: '/mongodb/docs', '/vercel/next.js')
 * @param topic ドキュメントのトピック (例: 'hooks', 'routing')
 * @param mode ドキュメントモード ('code' または 'info')
 * @param page ページ番号 (1から開始)
 * @returns ドキュメント内容
 */
export async function getLibraryDocs(
  libraryId: string,
  topic?: string,
  mode: 'code' | 'info' = 'code',
  page: number = 1
): Promise<Documentation> {
  console.log(`\n📖 Fetching documentation for: ${libraryId}`);
  if (topic) console.log(`   Topic: ${topic}`);
  console.log(`   Mode: ${mode}, Page: ${page}`);

  const args: any = {
    context7CompatibleLibraryID: libraryId,
    mode,
    page
  };

  if (topic) {
    args.topic = topic;
  }

  const result = await callMCPTool<any>('get-library-docs', args);

  console.log(`✅ Documentation retrieved successfully`);

  return {
    content: typeof result === 'string' ? result : JSON.stringify(result, null, 2),
    page,
    hasMore: true // Context7は最大10ページまでサポート
  };
}

/**
 * ライブラリ名からドキュメントを取得する完全なワークフロー
 * @param libraryName ライブラリ名
 * @param topic オプションのトピック
 * @param mode ドキュメントモード
 * @param page ページ番号
 * @returns ドキュメント内容
 */
export async function getDocsByLibraryName(
  libraryName: string,
  topic?: string,
  mode: 'code' | 'info' = 'code',
  page: number = 1
): Promise<Documentation> {
  try {
    // Step 1: ライブラリIDを解決
    const resolveResult = await resolveLibraryId(libraryName);

    if (!resolveResult.selected && (!resolveResult.matches || resolveResult.matches.length === 0)) {
      throw new Error(`No library found for: ${libraryName}`);
    }

    // 最も関連性の高いライブラリを選択
    const selectedLibrary = resolveResult.selected || resolveResult.matches[0];

    // Step 2: ドキュメントを取得
    const docs = await getLibraryDocs(selectedLibrary.id, topic, mode, page);

    return docs;
  } catch (error) {
    console.error(`\n❌ Error in workflow:`, error);
    throw error;
  }
}

/**
 * 複数ページのドキュメントを取得
 * @param libraryId ライブラリID
 * @param topic トピック
 * @param mode モード
 * @param maxPages 最大ページ数 (最大10)
 * @returns 全ページのドキュメント
 */
export async function getMultiplePages(
  libraryId: string,
  topic: string,
  mode: 'code' | 'info' = 'code',
  maxPages: number = 3
): Promise<Documentation[]> {
  const results: Documentation[] = [];
  const pageLimit = Math.min(maxPages, 10); // Context7の制限

  console.log(`\n📚 Fetching ${pageLimit} pages of documentation...`);

  for (let page = 1; page <= pageLimit; page++) {
    const docs = await getLibraryDocs(libraryId, topic, mode, page);
    results.push(docs);
  }

  console.log(`✅ Retrieved ${results.length} pages`);
  return results;
}

// メイン実行例
if (import.meta.url === `file://${process.argv[1]}`) {
  const args = process.argv.slice(2);

  if (args.length === 0) {
    console.log(`
Usage:
  node context7-workflow.js <library-name> [topic] [mode] [page]

Examples:
  node context7-workflow.js react
  node context7-workflow.js react hooks
  node context7-workflow.js react hooks code 1
  node context7-workflow.js next.js routing info 1

Modes:
  code - API references and code examples (default)
  info - Conceptual guides and architectural information
    `);
    process.exit(1);
  }

  const [libraryName, topic, mode = 'code', pageStr = '1'] = args;
  const page = parseInt(pageStr, 10);

  (async () => {
    try {
      const docs = await getDocsByLibraryName(
        libraryName,
        topic,
        mode as 'code' | 'info',
        page
      );

      console.log('\n' + '='.repeat(80));
      console.log('DOCUMENTATION:');
      console.log('='.repeat(80));
      console.log(docs.content);
      console.log('='.repeat(80));

    } catch (error) {
      console.error('Error:', error);
      process.exit(1);
    } finally {
      await closeConnection();
    }
  })();
}
