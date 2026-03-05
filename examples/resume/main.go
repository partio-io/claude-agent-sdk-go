// Example: resume a previous session by ID.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	claude "github.com/partio-io/claude-agent-sdk-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <session-id>\n", os.Args[0])
		os.Exit(1)
	}
	sessionID := os.Args[1]

	ctx := context.Background()

	session := claude.ResumeSession(sessionID,
		claude.WithModel("claude-sonnet-4-6"),
	)
	defer func() { _ = session.Close() }()

	if err := session.Send(ctx, "What number did I ask you to remember?"); err != nil {
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
