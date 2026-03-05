package process

import (
	"slices"
	"testing"
)

func TestBuildArgsMinimal(t *testing.T) {
	args := BuildArgs(Config{}, false)
	// Should always have --print and --output-format
	if !slices.Contains(args, "--print") {
		t.Error("missing --print")
	}
	if !slices.Contains(args, "stream-json") {
		t.Error("missing stream-json output format")
	}
}

func TestBuildArgsStreaming(t *testing.T) {
	args := BuildArgs(Config{}, true)
	if !slices.Contains(args, "--input-format") {
		t.Error("missing --input-format for streaming mode")
	}
}

func TestBuildArgsModel(t *testing.T) {
	args := BuildArgs(Config{Model: "claude-opus-4-6"}, false)
	idx := slices.Index(args, "--model")
	if idx < 0 {
		t.Fatal("missing --model")
	}
	if args[idx+1] != "claude-opus-4-6" {
		t.Errorf("model = %q", args[idx+1])
	}
}

func TestBuildArgsAllOptions(t *testing.T) {
	cfg := Config{
		Model:                  "claude-sonnet-4-6",
		SystemPrompt:           "Be helpful",
		AppendSystemPrompt:     "Extra instructions",
		Verbose:                true,
		AllowedTools:           []string{"Read", "Write"},
		DisallowedTools:        []string{"Bash"},
		PermissionMode:         "bypassPermissions",
		MaxTurns:               5,
		MaxBudgetUSD:           2.5,
		IncludePartialMessages: true,
		Resume:                 "session-123",
		ForkSession:            true,
		ContinueConversation:   false,
		NoSessionPersistence:   true,
		AddDirs:                []string{"/extra"},
	}
	args := BuildArgs(cfg, true)

	expects := []string{
		"--model", "claude-sonnet-4-6",
		"--system-prompt", "Be helpful",
		"--append-system-prompt", "Extra instructions",
		"--verbose",
		"--include-partial-messages",
		"--permission-mode", "bypassPermissions",
		"--max-turns", "5",
		"--resume", "session-123",
		"--fork-session",
		"--no-session-persistence",
		"--add-dir", "/extra",
	}

	for _, e := range expects {
		if !slices.Contains(args, e) {
			t.Errorf("missing expected arg: %s", e)
		}
	}
}

func TestBuildArgsNoZeroValues(t *testing.T) {
	args := BuildArgs(Config{}, false)
	// Zero values should not generate flags.
	for _, flag := range []string{"--model", "--system-prompt", "--max-turns", "--resume"} {
		if slices.Contains(args, flag) {
			t.Errorf("should not include %s with zero value", flag)
		}
	}
}
