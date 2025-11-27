import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { SSEClientTransport } from '@modelcontextprotocol/sdk/client/sse.js';
import { readFileSync } from 'fs';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

let cachedClient: Client | null = null;

interface MCPConfig {
  type: string;
  url?: string;
  command?: string;
  args?: string[];
  env?: Record<string, string>;
  headers?: Record<string, string>;
}

interface MCPServersConfig {
  mcpServers: Record<string, MCPConfig>;
}

function loadConfig(serverName: string): MCPConfig {
  // .mcp.jsonのパスを解決（pluginsディレクトリのルートから）
  const mcpConfigPath = resolve(__dirname, '../../../.mcp.json');

  try {
    const configContent = readFileSync(mcpConfigPath, 'utf-8');
    const config: MCPServersConfig = JSON.parse(configContent);

    if (!config.mcpServers || !config.mcpServers[serverName]) {
      throw new Error(`MCP server "${serverName}" not found in .mcp.json`);
    }

    return config.mcpServers[serverName];
  } catch (error) {
    throw new Error(`Failed to load MCP config: ${error}`);
  }
}

function extractPayload(result: any): any {
  // MCPツールの結果からペイロードを抽出
  if (result && result.content && Array.isArray(result.content)) {
    // contentが配列の場合、最初の要素のtextを返す
    const firstContent = result.content[0];
    if (firstContent && firstContent.type === 'text') {
      try {
        // JSONとしてパース可能な場合はパース
        return JSON.parse(firstContent.text);
      } catch {
        // パースできない場合はそのまま返す
        return firstContent.text;
      }
    }
  }

  // それ以外の場合はそのまま返す
  return result;
}

async function getClient(): Promise<Client> {
  if (cachedClient) {
    return cachedClient;
  }

  // .mcp.jsonから設定を読み込む
  const config = loadConfig('context7');

  if (config.type !== 'http') {
    throw new Error(`Unsupported MCP server type: ${config.type}`);
  }

  if (!config.url) {
    throw new Error('URL is required for HTTP MCP server');
  }

  const transport = new SSEClientTransport(
    new URL(config.url),
    {
      headers: config.headers || {}
    }
  );

  const client = new Client({
    name: 'context7-skill-client',
    version: '0.1.0'
  });

  await client.connect(transport);
  cachedClient = client;
  return client;
}

export async function callMCPTool<T>(toolName: string, input: any): Promise<T> {
  const client = await getClient();
  const result = await client.callTool({
    name: toolName,
    arguments: input
  });

  return extractPayload(result) as T;
}

export async function closeConnection() {
  if (cachedClient) {
    await cachedClient.close();
    cachedClient = null;
  }
}
