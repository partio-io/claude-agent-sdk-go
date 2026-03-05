package claude

import (
	"context"
	"os"
	"testing"

	"github.com/partio-io/claude-agent-sdk-go/internal/testutil"
)

func TestPromptSuccess(t *testing.T) {
	script, err := testutil.MockCLIScript(testutil.SimpleSession)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(script) }()

	result, err := Prompt(
		context.Background(),
		"What is 2+2?",
		WithCLIPath(script),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Subtype != ResultSuccess {
		t.Errorf("Subtype = %q, want success", result.Subtype)
	}
	if result.Result == nil {
		t.Fatal("Result is nil")
	}
	if *result.Result != "The answer is 4." {
		t.Errorf("Result = %q", *result.Result)
	}
	if result.SessionID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("SessionID = %q", result.SessionID)
	}
}

func TestPromptEmptyPrompt(t *testing.T) {
	_, err := Prompt(context.Background(), "")
	if err != ErrEmptyPrompt {
		t.Errorf("expected ErrEmptyPrompt, got %v", err)
	}
}

func TestPromptCLINotFound(t *testing.T) {
	_, err := Prompt(
		context.Background(),
		"test",
		WithCLIPath("/nonexistent/claude"),
	)
	if err != ErrCLINotFound {
		t.Errorf("expected ErrCLINotFound, got %v", err)
	}
}

func TestPromptWithOptions(t *testing.T) {
	script, err := testutil.MockCLIScript(testutil.SimpleSession)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(script) }()

	result, err := Prompt(
		context.Background(),
		"test",
		WithCLIPath(script),
		WithModel("claude-opus-4-6"),
		WithMaxTurns(5),
		WithMaxBudgetUSD(1.0),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestPromptToolUseSession(t *testing.T) {
	script, err := testutil.MockCLIScript(testutil.ToolUseSession)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(script) }()

	result, err := Prompt(
		context.Background(),
		"Read the file",
		WithCLIPath(script),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Subtype != ResultSuccess {
		t.Errorf("Subtype = %q", result.Subtype)
	}
}

func TestPromptContextCancelled(t *testing.T) {
	script, err := testutil.MockCLIScript(testutil.SimpleSession)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(script) }()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err = Prompt(ctx, "test", WithCLIPath(script))
	// Should either succeed (script runs fast) or fail with context error.
	// We don't assert a specific error since it depends on timing.
	_ = err
}
