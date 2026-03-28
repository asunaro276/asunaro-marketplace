package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	gitClient := newGitClient()
	githubClient := newGitHubClient(token)
	svc := newService(gitClient, githubClient)

	s := buildServer(svc)
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("MCP サーバーエラー: %v", err)
	}
}
