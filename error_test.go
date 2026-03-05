package claude

import (
	"errors"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"CLINotFound", ErrCLINotFound, "claude: CLI not found"},
		{"SessionClosed", ErrSessionClosed, "claude: session closed"},
		{"EmptyPrompt", ErrEmptyPrompt, "claude: empty prompt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessError(t *testing.T) {
	inner := errors.New("signal: killed")
	pe := &ProcessError{ExitCode: 137, Stderr: "OOM", Err: inner}

	if got := pe.Error(); got != "claude: process exited with code 137: OOM" {
		t.Errorf("unexpected error: %s", got)
	}

	if !errors.Is(pe, inner) {
		t.Error("ProcessError should wrap inner error")
	}
}

func TestProcessErrorNoStderr(t *testing.T) {
	pe := &ProcessError{ExitCode: 1}
	if got := pe.Error(); got != "claude: process exited with code 1" {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestParseError(t *testing.T) {
	inner := errors.New("unexpected EOF")
	pe := &ParseError{Raw: []byte(`{broken`), Err: inner}

	want := `claude: parse error: unexpected EOF (raw: {broken)`
	if got := pe.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if !errors.Is(pe, inner) {
		t.Error("ParseError should wrap inner error")
	}
}

func TestPermissionDeniedError(t *testing.T) {
	pe := &PermissionDeniedError{ToolName: "Bash", Reason: "not allowed"}
	want := `claude: permission denied for tool "Bash": not allowed`
	if got := pe.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestVersionError(t *testing.T) {
	ve := &VersionError{Got: "1.0", Want: "2.0"}
	want := `claude: version mismatch: got 1.0, want 2.0`
	if got := ve.Error(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
