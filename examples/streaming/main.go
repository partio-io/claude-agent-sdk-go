// Example: token-level streaming with partial messages.
package main

import (
	"context"
	"fmt"
	"log"

	claude "github.com/anthropics/claude-agent-sdk-go"
)

func main() {
	ctx := context.Background()

	session := claude.NewSession(
		claude.WithModel("claude-sonnet-4-6"),
		claude.WithIncludePartialMessages(true),
	)
	defer func() { _ = session.Close() }()

	if err := session.Send(ctx, "Write a haiku about Go programming"); err != nil {
		log.Fatal(err)
	}

	for msg, err := range session.Stream(ctx) {
		if err != nil {
			log.Fatal(err)
		}
		switch m := msg.(type) {
		case *claude.StreamEvent:
			// Token-level streaming events.
			if delta, ok := m.Event["delta"].(map[string]any); ok {
				if text, ok := delta["text"].(string); ok {
					fmt.Print(text)
				}
			}
		case *claude.AssistantMessage:
			fmt.Println() // newline after streaming
		case *claude.ResultMessage:
			fmt.Printf("\n(tokens: in=%d out=%d)\n",
				m.Usage.InputTokens, m.Usage.OutputTokens)
		}
	}
}
