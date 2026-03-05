package testutil

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// MockCLIScript creates a temporary shell script that outputs NDJSON lines.
// Returns the path to the script. Caller must remove it when done.
func MockCLIScript(lines []string) (string, error) {
	var sb strings.Builder
	sb.WriteString("#!/bin/sh\n")

	for _, line := range lines {
		// Escape single quotes in JSON.
		escaped := strings.ReplaceAll(line, "'", "'\\''")
		fmt.Fprintf(&sb, "echo '%s'\n", escaped)
	}

	f, err := os.CreateTemp("", "mock-claude-*.sh")
	if err != nil {
		return "", err
	}
	if _, err := f.WriteString(sb.String()); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return "", err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(f.Name())
		return "", err
	}

	if err := os.Chmod(f.Name(), 0o755); err != nil {
		_ = os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

// MockStreamingCLIScript creates a script that reads stdin and outputs NDJSON.
// It echoes the initial lines, then for each stdin line reads, echoes the response lines.
func MockStreamingCLIScript(initLines []string, responseLines []string) (string, error) {
	var sb strings.Builder
	sb.WriteString("#!/bin/sh\n")

	// Output initial lines (system init).
	for _, line := range initLines {
		escaped := strings.ReplaceAll(line, "'", "'\\''")
		fmt.Fprintf(&sb, "echo '%s'\n", escaped)
	}

	// Read from stdin and output response for each user message.
	sb.WriteString("while IFS= read -r input; do\n")
	// Skip control messages (type: control_response)
	sb.WriteString("  case \"$input\" in *control_response*) continue ;; esac\n")
	for _, line := range responseLines {
		escaped := strings.ReplaceAll(line, "'", "'\\''")
		fmt.Fprintf(&sb, "  echo '%s'\n", escaped)
	}
	sb.WriteString("done\n")

	f, err := os.CreateTemp("", "mock-claude-stream-*.sh")
	if err != nil {
		return "", err
	}
	if _, err := f.WriteString(sb.String()); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return "", err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(f.Name())
		return "", err
	}

	if err := os.Chmod(f.Name(), 0o755); err != nil {
		_ = os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

// IsHelperProcess returns true if the current process is a test helper subprocess.
func IsHelperProcess() bool {
	return os.Getenv("GO_TEST_HELPER_PROCESS") == "1"
}

// HelperProcess runs the mock CLI helper if GO_TEST_HELPER_PROCESS is set.
// Call this from TestMain or at the start of tests that use exec.Command with the test binary.
func HelperProcess() {
	if !IsHelperProcess() {
		return
	}
	lines := os.Getenv("GO_TEST_NDJSON_LINES")
	if lines == "" {
		os.Exit(0)
	}
	var ndjson []string
	if err := json.Unmarshal([]byte(lines), &ndjson); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse NDJSON lines: %v\n", err)
		os.Exit(1)
	}
	for _, line := range ndjson {
		fmt.Println(line)
	}
	os.Exit(0)
}

// HelperCommand returns an exec.Cmd that runs the test binary as a helper process.
func HelperCommand(testBinary string, lines []string) *exec.Cmd {
	linesJSON, _ := json.Marshal(lines) //nolint:errcheck // lines is always []string
	cmd := exec.Command(testBinary)
	cmd.Env = append(os.Environ(),
		"GO_TEST_HELPER_PROCESS=1",
		"GO_TEST_NDJSON_LINES="+string(linesJSON),
	)
	return cmd
}
