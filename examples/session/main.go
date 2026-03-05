// Example: multi-turn session using claude.NewSession.
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
		claude.WithMaxTurns(3),
	)
	defer func() { _ = session.Close() }()

	// First turn
	if err := session.Send(ctx, "What is 5 + 3?"); err != nil {
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
					fmt.Printf("[Turn 1] %s\n", tb.Text)
				}
			}
		case *claude.ResultMessage:
			fmt.Printf("[Turn 1] Done in %d turns\n", m.NumTurns)
		}
	}

	// Second turn (Claude remembers context)
	if err := session.Send(ctx, "Multiply that by 2"); err != nil {
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
					fmt.Printf("[Turn 2] %s\n", tb.Text)
				}
			}
		case *claude.ResultMessage:
			fmt.Printf("[Turn 2] Done in %d turns\n", m.NumTurns)
		}
	}

	fmt.Printf("Session ID: %s\n", session.SessionID())
}
