package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// buildServer は mcp-go を使って MCP サーバーを構築する。
// MCP ライブラリを差し替えたい場合はこのファイルのみ修正する。
func buildServer(svc *Service) *server.MCPServer {
	s := server.NewMCPServer("pr-context-mcp", "1.0.0",
		server.WithToolCapabilities(true),
	)

	tool := mcp.NewTool("get_pr_context",
		mcp.WithDescription(
			"ファイルの特定行に関連する PR のコンテキスト（ディスクリプション・コメント）を取得する。"+
				"コードがなぜそのように書かれているかを調査するのに役立つ。",
		),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("調査するファイルのパス"),
		),
		mcp.WithNumber("line_number",
			mcp.Required(),
			mcp.Description("調査する行番号（1始まり）"),
		),
		mcp.WithString("repo",
			mcp.Description("リポジトリ（owner/repo 形式）。省略時は git remote から自動検出"),
		),
		mcp.WithString("cwd",
			mcp.Description("git リポジトリのルートディレクトリ。省略時はカレントディレクトリ"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, err := req.RequireString("file_path")
		if err != nil {
			return mcp.NewToolResultError("file_path は必須です"), nil
		}

		lineNumber := req.GetInt("line_number", 0)
		if lineNumber <= 0 {
			return mcp.NewToolResultError("line_number は1以上の整数を指定してください"), nil
		}

		repo := req.GetString("repo", "")
		cwd := req.GetString("cwd", "")

		result, err := svc.GetPRContext(filePath, lineNumber, repo, cwd)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	return s
}
