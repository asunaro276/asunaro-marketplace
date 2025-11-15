"""
Figma MCP: get_design_context

レイヤーまたは選択範囲のデザインコンテキストを取得します。
デフォルトでは React + Tailwind CSS で出力されますが、
フレームワークやコンポーネントライブラリをカスタマイズ可能です。

使用方法:
    Claude が直接 MCP ツールを呼び出します：
    mcp__figma_dev_mode_mcp__get_design_context()

戻り値の例:
    {
        "component": "Button",
        "code": "<button className=\"bg-blue-500 ...\">\n  Click me\n</button>",
        "framework": "react",
        "styling": "tailwind"
    }
"""

# このファイルは Claude がツールの仕様を理解するためのリファレンスです。
# 実際の実装は Figma Dev Mode MCP サーバーが提供します。


def get_design_context(
    framework: str = "react",
    styling: str = "tailwind",
    custom_instructions: str = ""
) -> dict:
    """
    選択された Figma レイヤーのデザインコンテキストをコードとして取得します。

    Parameters:
        framework: フレームワーク指定（react, vue, html など）
        styling: スタイリング方法（tailwind, css, scss など）
        custom_instructions: カスタム指示（例: "shadcn/ui を使用"）

    Returns:
        デザインコンテキスト情報を含む辞書

    Notes:
        - Vue を使用する場合は framework="vue" を指定
        - デフォルトは React + Tailwind CSS
        - レイアウトの忠実度を維持するため、get_screenshot と併用推奨
    """
    # MCP ツールの呼び出し例:
    # result = call_mcp_tool("mcp__figma_dev_mode_mcp__get_design_context", {
    #     "framework": framework,
    #     "styling": styling,
    #     "custom_instructions": custom_instructions
    # })
    # return result
    pass
