// Example: one-shot prompt using claude.Prompt.
package main

import (
	"context"
	"fmt"
	"log"

	claude "github.com/partio-io/claude-agent-sdk-go"
)

func main() {
	ctx := context.Background()

	result, err := claude.Prompt(ctx, "What is 2+2?",
		claude.WithModel("claude-sonnet-4-6"),
		claude.WithMaxTurns(1),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Subtype: %s\n", result.Subtype)
	if result.Result != nil {
		fmt.Printf("Result: %s\n", *result.Result)
	}
	if result.TotalCostUSD != nil {
		fmt.Printf("Cost: $%.4f\n", *result.TotalCostUSD)
	}
}
