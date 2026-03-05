// Example: MCP server integration.
package main

import (
	"context"
	"fmt"
	"log"

	claude "github.com/partio-io/claude-agent-sdk-go"
)

func main() {
	ctx := context.Background()

	session := claude.NewSession(
		claude.WithModel("claude-sonnet-4-6"),
		claude.WithMCPServer("filesystem", &claude.MCPStdioServer{
			Command: "npx",
			Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
		}),
	)
	defer func() { _ = session.Close() }()

	if err := session.Send(ctx, "List the files in /tmp using the filesystem MCP server"); err != nil {
		log.Fatal(err)
	}

	for msg, err := range session.Stream(ctx) {
		if err != nil {
			log.Fatal(err)
		}
		switch m := msg.(type) {
		case *claude.AssistantMessage:
			for _, block := range m.Message.Content {
				if tb, ok := block.(*claude.TextBlock); ok {
					fmt.Println(tb.Text)
				}
			}
		case *claude.ResultMessage:
			fmt.Printf("Done (cost: $%.4f)\n", safeFloat(m.TotalCostUSD))
		}
	}
}

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
