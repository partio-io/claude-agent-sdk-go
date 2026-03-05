// Example: hooks to intercept and control tool use.
package main

import (
	"context"
	"fmt"
	"log"

	claude "github.com/anthropics/claude-agent-sdk-go"
)

func main() {
	ctx := context.Background()

	matcher := "Bash|Edit|Write"
	session := claude.NewSession(
		claude.WithModel("claude-sonnet-4-6"),
		claude.WithHook(claude.HookPreToolUse, claude.HookMatcher{
			Matcher: &matcher,
			Handler: func(ctx context.Context, input claude.HookCallbackInput) (claude.HookOutput, error) {
				fmt.Printf("[Hook] Pre-tool: %s\n", input.ToolName)

				// Auto-approve read-only tools.
				if input.ToolName == "Read" || input.ToolName == "Glob" || input.ToolName == "Grep" {
					return claude.HookOutput{Decision: "allow"}, nil
				}

				// Allow all others for this example.
				fmt.Printf("[Hook] Allowing %s\n", input.ToolName)
				return claude.HookOutput{Decision: "allow"}, nil
			},
		}),
		claude.WithHook(claude.HookPostToolUse, claude.HookMatcher{
			Handler: func(ctx context.Context, input claude.HookCallbackInput) (claude.HookOutput, error) {
				fmt.Printf("[Hook] Post-tool: %s completed\n", input.ToolName)
				return claude.HookOutput{}, nil
			},
		}),
	)
	defer func() { _ = session.Close() }()

	if err := session.Send(ctx, "List the files in the current directory"); err != nil {
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
			fmt.Printf("Done (turns: %d)\n", m.NumTurns)
		}
	}
}
