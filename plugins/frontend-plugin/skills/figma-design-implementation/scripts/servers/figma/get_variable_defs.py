"""
Figma MCP: get_variable_defs

Figma 選択範囲で使用されている変数とスタイルを取得します。
色、間隔、タイポグラフィなどのデザイントークンを抽出します。

使用方法:
    Claude が直接 MCP ツールを呼び出します：
    mcp__figma_dev_mode_mcp__get_variable_defs()

戻り値の例:
    {
        "colors": {
            "primary": "#3B82F6",
            "secondary": "#10B981",
            "text-primary": "#1F2937"
        },
        "spacing": {
            "sm": "8px",
            "md": "16px",
            "lg": "24px"
        },
        "typography": {
            "heading-1": {
                "fontFamily": "Inter",
                "fontSize": "32px",
                "fontWeight": "700",
                "lineHeight": "1.2"
            }
        }
    }
"""


def get_variable_defs() -> dict:
    """
    選択された Figma 要素で使用されているデザイントークンを取得します。

    Returns:
        デザイントークンを含む辞書
        - colors: カラーパレット
        - spacing: スペーシング値
        - typography: タイポグラフィ定義
        - effects: エフェクト（シャドウ、ブラーなど）

    Notes:
        - このツールで取得したトークンは CSS カスタムプロパティに変換可能
        - design_tokens_to_css() ユーティリティと併用推奨
    """
    # MCP ツールの呼び出し例:
    # result = call_mcp_tool("mcp__figma_dev_mode_mcp__get_variable_defs", {})
    # return result
    pass
