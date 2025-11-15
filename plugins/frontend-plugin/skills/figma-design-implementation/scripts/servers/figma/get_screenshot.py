"""
Figma MCP: get_screenshot

選択範囲のスクリーンショットを取得します。
レイアウトの忠実度を維持する際に推奨されます。

使用方法:
    Claude が直接 MCP ツールを呼び出します：
    mcp__figma_dev_mode_mcp__get_screenshot()

戻り値の例:
    {
        "image_url": "data:image/png;base64,...",
        "width": 1200,
        "height": 800,
        "format": "png"
    }
"""


def get_screenshot(scale: float = 2.0, format: str = "png") -> dict:
    """
    選択された Figma 要素のスクリーンショットを取得します。

    Parameters:
        scale: スケール倍率（デフォルト: 2.0 for Retina）
        format: 画像フォーマット（png, jpg）

    Returns:
        スクリーンショット情報を含む辞書
        - image_url: Base64エンコードされた画像データまたはURL
        - width: 幅（ピクセル）
        - height: 高さ（ピクセル）

    Notes:
        - get_design_context と併用することで、視覚的な忠実度を確保
        - 複雑なレイアウトやグラデーションの再現に有効
    """
    # MCP ツールの呼び出し例:
    # result = call_mcp_tool("mcp__figma_dev_mode_mcp__get_screenshot", {
    #     "scale": scale,
    #     "format": format
    # })
    # return result
    pass
