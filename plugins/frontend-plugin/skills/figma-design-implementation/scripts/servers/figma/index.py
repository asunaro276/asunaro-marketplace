"""
Figma MCP サーバー - インデックス

このモジュールは Figma Dev Mode MCP サーバーが提供する
すべてのツールへのエントリーポイントです。

利用可能なツール:
    - get_design_context: デザインをコードに変換
    - get_variable_defs: デザイントークンを取得
    - get_screenshot: スクリーンショットを取得
    - get_code_connect_map: デザインとコードのマッピング
    - create_design_system_rules: デザインシステムルールを生成

使用方法:
    Claude はこのディレクトリ構造を探索し、必要なツールファイルを
    読み込んで仕様を理解します。実際のツール呼び出しは MCP プロトコル
    経由で行われます。
"""

from . import get_design_context
from . import get_variable_defs
from . import get_screenshot
from . import get_code_connect_map
from . import create_design_system_rules

__all__ = [
    'get_design_context',
    'get_variable_defs',
    'get_screenshot',
    'get_code_connect_map',
    'create_design_system_rules',
]


# MCP サーバー接続情報
MCP_SERVER_CONFIG = {
    "name": "Figma Dev Mode MCP",
    "type": "sse",
    "url": "http://127.0.0.1:3845/sse",
    "description": "Figma デザインファイルからデザイン情報を取得し、コードに変換"
}


def list_available_tools() -> list[str]:
    """
    利用可能な Figma MCP ツールのリストを返します。

    Returns:
        ツール名のリスト
    """
    return [
        "get_design_context",
        "get_variable_defs",
        "get_screenshot",
        "get_code_connect_map",
        "create_design_system_rules"
    ]


def get_tool_info(tool_name: str) -> dict:
    """
    指定されたツールの情報を取得します。

    Parameters:
        tool_name: ツール名

    Returns:
        ツール情報を含む辞書
    """
    tool_docs = {
        "get_design_context": {
            "description": "レイヤーまたは選択範囲のデザインコンテキストをコードとして取得",
            "use_case": "Figma デザインを HTML/CSS/Vue に変換",
            "file": "get_design_context.py"
        },
        "get_variable_defs": {
            "description": "デザイントークン（色、間隔、タイポグラフィ）を取得",
            "use_case": "CSS カスタムプロパティの生成",
            "file": "get_variable_defs.py"
        },
        "get_screenshot": {
            "description": "選択範囲のスクリーンショットを取得",
            "use_case": "視覚的な忠実度の確保",
            "file": "get_screenshot.py"
        },
        "get_code_connect_map": {
            "description": "Figma とコードのマッピングを取得",
            "use_case": "既存コンポーネントの再利用",
            "file": "get_code_connect_map.py"
        },
        "create_design_system_rules": {
            "description": "デザインシステムルールを生成",
            "use_case": "デザインシステムの文脈理解",
            "file": "create_design_system_rules.py"
        }
    }

    return tool_docs.get(tool_name, {"error": "ツールが見つかりません"})


if __name__ == "__main__":
    print("=== Figma MCP サーバー ツール一覧 ===")
    for tool in list_available_tools():
        info = get_tool_info(tool)
        print(f"\n{tool}:")
        print(f"  説明: {info['description']}")
        print(f"  用途: {info['use_case']}")
