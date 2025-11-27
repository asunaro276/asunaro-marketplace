#!/usr/bin/env python3
"""
Context7 MCPクライアント サンプルスクリプト

このスクリプトは、Context7 MCPサーバーを使用してライブラリのドキュメントを
取得するためのテンプレートです。プロジェクトの要件に応じてカスタマイズしてください。

使用方法:
    python context7_client_example.py

依存関係:
    pip install mcp anthropic

参考:
    - Context7: https://github.com/upstash/context7
    - MCP Best Practices: https://www.anthropic.com/engineering/code-execution-with-mcp
"""

import asyncio
import json
from typing import Optional, Dict, Any

try:
    from mcp import ClientSession, StdioServerParameters
    from mcp.client.stdio import stdio_client
except ImportError:
    print("エラー: MCPライブラリがインストールされていません")
    print("インストールするには: pip install mcp")
    exit(1)


class Context7Client:
    """Context7 MCPサーバーとの通信を管理するクライアント"""

    def __init__(self):
        self.session: Optional[ClientSession] = None
        self.server_params = StdioServerParameters(
            command="npx",
            args=["-y", "@upstash/context7"]
        )

    async def __aenter__(self):
        """非同期コンテキストマネージャーのエントリー"""
        # MCPサーバーへの接続を確立
        self.stdio_transport = await stdio_client(self.server_params).__aenter__()
        self.session = self.stdio_transport[1]

        # サーバーの初期化
        await self.session.initialize()

        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """非同期コンテキストマネージャーの終了"""
        if self.stdio_transport:
            await self.stdio_transport[0].__aexit__(exc_type, exc_val, exc_tb)

    async def resolve_library_id(self, library_name: str) -> Dict[str, Any]:
        """
        ライブラリ名をContext7互換のライブラリIDに解決する

        Args:
            library_name: 検索するライブラリ名（例: "react", "next.js"）

        Returns:
            解決されたライブラリ情報を含む辞書

        Example:
            result = await client.resolve_library_id("react")
            library_id = result["library_id"]  # "/facebook/react"
        """
        if not self.session:
            raise RuntimeError("セッションが初期化されていません")

        try:
            response = await self.session.call_tool(
                "resolve-library-id",
                arguments={"libraryName": library_name}
            )

            # レスポンスを解析
            if response.content:
                # MCPツールのレスポンスは通常テキスト形式
                result_text = response.content[0].text if response.content else ""

                print(f"✓ ライブラリID解決成功: {library_name}")
                print(f"  レスポンス: {result_text}")

                return {
                    "library_name": library_name,
                    "response": result_text,
                    "success": True
                }
            else:
                print(f"✗ ライブラリIDの解決に失敗: {library_name}")
                return {"success": False, "error": "レスポンスが空です"}

        except Exception as e:
            print(f"✗ エラーが発生しました: {e}")
            return {"success": False, "error": str(e)}

    async def get_library_docs(
        self,
        library_id: str,
        topic: Optional[str] = None,
        mode: str = "code",
        page: int = 1
    ) -> Dict[str, Any]:
        """
        指定されたライブラリのドキュメントを取得する

        Args:
            library_id: Context7互換のライブラリID（例: "/facebook/react"）
            topic: フォーカスするトピック（例: "hooks", "routing"）
            mode: ドキュメントモード（"code" または "info"）
            page: ページ番号（1-10）

        Returns:
            ドキュメント情報を含む辞書

        Example:
            docs = await client.get_library_docs(
                library_id="/facebook/react",
                topic="hooks",
                mode="code",
                page=1
            )
        """
        if not self.session:
            raise RuntimeError("セッションが初期化されていません")

        # パラメータの構築
        arguments = {
            "context7CompatibleLibraryID": library_id,
            "mode": mode,
            "page": page
        }

        if topic:
            arguments["topic"] = topic

        try:
            response = await self.session.call_tool(
                "get-library-docs",
                arguments=arguments
            )

            # レスポンスを解析
            if response.content:
                result_text = response.content[0].text if response.content else ""

                print(f"✓ ドキュメント取得成功: {library_id}")
                if topic:
                    print(f"  トピック: {topic}")
                print(f"  モード: {mode}, ページ: {page}")
                print(f"\n--- ドキュメント内容 ---")
                print(result_text)
                print(f"--- 終了 ---\n")

                return {
                    "library_id": library_id,
                    "topic": topic,
                    "mode": mode,
                    "page": page,
                    "content": result_text,
                    "success": True
                }
            else:
                print(f"✗ ドキュメント取得に失敗: {library_id}")
                return {"success": False, "error": "レスポンスが空です"}

        except Exception as e:
            print(f"✗ エラーが発生しました: {e}")
            return {"success": False, "error": str(e)}


async def example_workflow():
    """
    Context7を使用した基本的なワークフローの例

    このワークフローをプロジェクトの要件に応じてカスタマイズしてください。
    """

    print("=" * 60)
    print("Context7 MCPクライアント サンプルワークフロー")
    print("=" * 60)
    print()

    async with Context7Client() as client:
        # === ステップ1: ライブラリIDの解決 ===
        print("ステップ1: ライブラリIDの解決")
        print("-" * 60)

        library_name = "react"  # ← ここをカスタマイズ
        resolve_result = await client.resolve_library_id(library_name)

        if not resolve_result.get("success"):
            print("ライブラリIDの解決に失敗しました。終了します。")
            return

        # NOTE: 実際には、resolve_library_idのレスポンスからライブラリIDを抽出する
        # ロジックを実装する必要があります。ここでは簡単のため、手動で指定します。
        library_id = "/facebook/react"  # ← 実際にはレスポンスから抽出

        print()

        # === ステップ2: ドキュメントの取得 ===
        print("ステップ2: ドキュメントの取得")
        print("-" * 60)

        docs_result = await client.get_library_docs(
            library_id=library_id,
            topic="hooks",  # ← ここをカスタマイズ
            mode="code",    # ← "code" または "info"
            page=1          # ← ページ番号（1-10）
        )

        if not docs_result.get("success"):
            print("ドキュメントの取得に失敗しました。")
            return

        print()

        # === ステップ3: 追加ページの取得（オプション） ===
        # コンテキストが不十分な場合は、追加ページを取得
        # print("ステップ3: 追加ページの取得")
        # print("-" * 60)
        #
        # docs_page2 = await client.get_library_docs(
        #     library_id=library_id,
        #     topic="hooks",
        #     mode="code",
        #     page=2
        # )
        #
        # print()

        print("=" * 60)
        print("ワークフロー完了")
        print("=" * 60)


async def custom_example():
    """
    カスタマイズ例: 複数のトピックについてドキュメントを取得

    この関数をテンプレートとして、プロジェクト固有のスクリプトを作成してください。
    """

    print("=" * 60)
    print("カスタム例: 複数トピックのドキュメント取得")
    print("=" * 60)
    print()

    # カスタマイズポイント
    library_name = "react"
    library_id = "/facebook/react"  # または resolve_library_id で取得
    topics = ["hooks", "context", "useEffect"]  # ← 取得したいトピックのリスト
    mode = "code"  # または "info"

    async with Context7Client() as client:
        for topic in topics:
            print(f"\n📚 トピック '{topic}' のドキュメントを取得中...")
            print("-" * 60)

            result = await client.get_library_docs(
                library_id=library_id,
                topic=topic,
                mode=mode,
                page=1
            )

            if result.get("success"):
                # ここで取得したドキュメントを使用
                # 例: ファイルに保存、データベースに格納、コード生成など
                content = result.get("content", "")

                # 保存例
                filename = f"docs_{topic.replace(' ', '_')}.md"
                with open(filename, "w", encoding="utf-8") as f:
                    f.write(f"# {topic}\n\n")
                    f.write(content)

                print(f"✓ ドキュメントを {filename} に保存しました")
            else:
                print(f"✗ '{topic}' のドキュメント取得に失敗")

            print()

    print("=" * 60)
    print("すべてのトピックの処理が完了しました")
    print("=" * 60)


def main():
    """
    メイン関数

    使用したいワークフローのコメントを外して実行してください。
    """

    # 基本的なワークフロー
    asyncio.run(example_workflow())

    # カスタム例（複数トピック）
    # asyncio.run(custom_example())


if __name__ == "__main__":
    main()
