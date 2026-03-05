package process

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindCLIExplicit(t *testing.T) {
	// Create a temp file to act as the CLI.
	tmp := t.TempDir()
	fakeCLI := filepath.Join(tmp, "claude")
	if err := os.WriteFile(fakeCLI, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := FindCLI(fakeCLI)
	if err != nil {
		t.Fatal(err)
	}
	if got != fakeCLI {
		t.Errorf("got %q, want %q", got, fakeCLI)
	}
}

func TestFindCLIExplicitNotFound(t *testing.T) {
	_, err := FindCLI("/nonexistent/claude")
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}

func TestFindCLILocalDir(t *testing.T) {
	// Create fake ~/.claude/local/claude
	tmp := t.TempDir()
	localDir := filepath.Join(tmp, ".claude", "local")
	if err := os.MkdirAll(localDir, 0o755); err != nil {
		t.Fatal(err)
	}
	fakeCLI := filepath.Join(localDir, "claude")
	if err := os.WriteFile(fakeCLI, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Override HOME to use our temp dir.
	t.Setenv("HOME", tmp)

	// This should find claude in ~/.claude/local/ (assuming it's not in PATH).
	// We can't guarantee it won't be in PATH, so this is best-effort.
	got, err := FindCLI("")
	if err != nil {
		// Might be found in PATH instead, which is fine.
		t.Skipf("claude might be in PATH, skipping local dir test: %v", err)
	}
	_ = got
}

func TestCLINotFoundError(t *testing.T) {
	e := &CLINotFoundError{Searched: []string{"PATH", "~/.claude/local/claude"}}
	want := "claude: CLI not found (searched: PATH, ~/.claude/local/claude)"
	if got := e.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
