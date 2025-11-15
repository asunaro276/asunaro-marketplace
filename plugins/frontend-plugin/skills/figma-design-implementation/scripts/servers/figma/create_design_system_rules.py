"""
Figma MCP: create_design_system_rules

デザインシステムの文脈を提供するルールファイルを生成します。
エージェントがデザインシステムの規約を理解するために使用します。

使用方法:
    Claude が直接 MCP ツールを呼び出します：
    mcp__figma_dev_mode_mcp__create_design_system_rules()

戻り値の例:
    {
        "rules": {
            "colors": {
                "naming_convention": "semantic",
                "usage_guidelines": "primary は CTA ボタンに使用"
            },
            "spacing": {
                "scale": "8px ベース",
                "usage": "コンポーネント間は 16px または 24px"
            },
            "components": {
                "button": {
                    "variants": ["primary", "secondary", "ghost"],
                    "sizes": ["sm", "md", "lg"]
                }
            }
        }
    }
"""


def create_design_system_rules() -> dict:
    """
    デザインシステムのルールと規約を生成します。

    Returns:
        デザインシステムルールを含む辞書
        - colors: カラーシステムの規約
        - spacing: スペーシングシステムの規約
        - components: コンポーネントの定義と使用方法

    Notes:
        - プロジェクト全体の一貫性を保つために重要
        - 生成されたルールは references/design-system-rules.md として保存推奨
        - 新しいコンポーネント実装時の参照資料として活用
    """
    # MCP ツールの呼び出し例:
    # result = call_mcp_tool("mcp__figma_dev_mode_mcp__create_design_system_rules", {})
    # return result
    pass
