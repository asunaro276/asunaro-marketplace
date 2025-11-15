"""
デザイントークンから CSS 変換ユーティリティ

Figma MCP の get_variable_defs で取得したデザイントークンを
CSS カスタムプロパティ（CSS変数）に変換します。
"""

from typing import Dict, Any


def design_tokens_to_css(tokens: Dict[str, Any], prefix: str = "") -> str:
    """
    デザイントークンを CSS カスタムプロパティに変換します。

    Args:
        tokens: get_variable_defs から取得したデザイントークン
        prefix: CSS変数のプレフィックス（オプション）

    Returns:
        CSS カスタムプロパティの文字列

    Example:
        >>> tokens = {
        ...     "colors": {
        ...         "primary": "#3B82F6",
        ...         "secondary": "#10B981"
        ...     },
        ...     "spacing": {
        ...         "sm": "8px",
        ...         "md": "16px"
        ...     }
        ... }
        >>> print(design_tokens_to_css(tokens))
        :root {
          --color-primary: #3B82F6;
          --color-secondary: #10B981;
          --spacing-sm: 8px;
          --spacing-md: 16px;
        }
    """
    css_vars = [":root {"]

    def process_tokens(obj: Dict[str, Any], current_prefix: str = ""):
        """再帰的にトークンを処理してCSS変数を生成"""
        for key, value in obj.items():
            # キーをケバブケースに変換
            kebab_key = key.replace("_", "-").replace(" ", "-").lower()
            variable_name = f"{current_prefix}-{kebab_key}" if current_prefix else kebab_key

            if isinstance(value, dict):
                # ネストされたオブジェクトの場合は再帰処理
                process_tokens(value, variable_name)
            else:
                # プリミティブ値の場合はCSS変数として追加
                full_name = f"{prefix}-{variable_name}" if prefix else variable_name
                css_vars.append(f"  --{full_name}: {value};")

    process_tokens(tokens)
    css_vars.append("}")

    return "\n".join(css_vars)


def tokens_to_tailwind_config(tokens: Dict[str, Any]) -> str:
    """
    デザイントークンを Tailwind CSS の設定形式に変換します。

    Args:
        tokens: get_variable_defs から取得したデザイントークン

    Returns:
        Tailwind CSS theme 設定の文字列（JavaScript）

    Example:
        >>> tokens = {"colors": {"primary": "#3B82F6"}}
        >>> print(tokens_to_tailwind_config(tokens))
        module.exports = {
          theme: {
            extend: {
              colors: {
                primary: '#3B82F6'
              }
            }
          }
        }
    """
    import json  # ここで必要なのでインポート

    # Tailwind の theme.extend 構造に変換
    theme_extend = {}

    if "colors" in tokens:
        theme_extend["colors"] = tokens["colors"]

    if "spacing" in tokens:
        theme_extend["spacing"] = tokens["spacing"]

    if "typography" in tokens:
        # タイポグラフィは fontSize と fontFamily に分離
        if "fontSize" not in theme_extend:
            theme_extend["fontSize"] = {}
        if "fontFamily" not in theme_extend:
            theme_extend["fontFamily"] = {}

        for key, value in tokens["typography"].items():
            if isinstance(value, dict):
                if "fontSize" in value:
                    theme_extend["fontSize"][key] = value["fontSize"]
                if "fontFamily" in value:
                    theme_extend["fontFamily"][key] = value["fontFamily"]

    # JavaScript 形式で出力
    config = f"""module.exports = {{
  theme: {{
    extend: {json.dumps(theme_extend, indent=6)}
  }}
}}"""

    return config


def tokens_to_scss_variables(tokens: Dict[str, Any]) -> str:
    """
    デザイントークンを SCSS 変数に変換します。

    Args:
        tokens: get_variable_defs から取得したデザイントークン

    Returns:
        SCSS 変数の文字列

    Example:
        >>> tokens = {"colors": {"primary": "#3B82F6"}}
        >>> print(tokens_to_scss_variables(tokens))
        $color-primary: #3B82F6;
    """
    scss_vars = []

    def process_tokens(obj: Dict[str, Any], prefix: str = ""):
        """再帰的にトークンを処理してSCSS変数を生成"""
        for key, value in obj.items():
            kebab_key = key.replace("_", "-").replace(" ", "-").lower()
            variable_name = f"{prefix}-{kebab_key}" if prefix else kebab_key

            if isinstance(value, dict):
                process_tokens(value, variable_name)
            else:
                scss_vars.append(f"${variable_name}: {value};")

    process_tokens(tokens)

    return "\n".join(scss_vars)


if __name__ == "__main__":
    # テスト例
    sample_tokens = {
        "colors": {
            "primary": "#3B82F6",
            "secondary": "#10B981",
            "text": {
                "primary": "#1F2937",
                "secondary": "#6B7280"
            }
        },
        "spacing": {
            "xs": "4px",
            "sm": "8px",
            "md": "16px",
            "lg": "24px"
        },
        "typography": {
            "heading-1": {
                "fontSize": "32px",
                "fontFamily": "Inter",
                "fontWeight": "700"
            }
        }
    }

    print("=== CSS カスタムプロパティ ===")
    print(design_tokens_to_css(sample_tokens))
    print("\n=== SCSS 変数 ===")
    print(tokens_to_scss_variables(sample_tokens))
    print("\n=== Tailwind Config ===")
    print(tokens_to_tailwind_config(sample_tokens))
