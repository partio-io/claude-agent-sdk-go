package claude

import (
	"context"
	"os"
	"testing"

	"github.com/partio-io/claude-agent-sdk-go/internal/testutil"
)

func TestNewSession(t *testing.T) {
	s := NewSession(WithModel("claude-sonnet-4-6"))
	defer func() { _ = s.Close() }()

	if s.cfg.model != "claude-sonnet-4-6" {
		t.Errorf("model = %q", s.cfg.model)
	}
}

func TestResumeSession(t *testing.T) {
	s := ResumeSession("session-abc", WithModel("claude-sonnet-4-6"))
	defer func() { _ = s.Close() }()

	if s.cfg.resume != "session-abc" {
		t.Errorf("resume = %q", s.cfg.resume)
	}
}

func TestSessionSendAndStream(t *testing.T) {
	responseLines := []string{
		testutil.AssistantTextOnly,
		testutil.ResultSuccess,
	}
	script, err := testutil.MockStreamingCLIScript(
		[]string{testutil.SystemInit},
		responseLines,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(script) }()

	session := NewSession(WithCLIPath(script))
	defer func() { _ = session.Close() }()

	ctx := context.Background()

	if err := session.Send(ctx, "What is 2+2?"); err != nil {
		t.Fatal(err)
	}

	var messages []Message
	for msg, err := range session.Stream(ctx) {
		if err != nil {
			t.Fatal(err)
		}
		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		t.Fatal("no messages received")
	}

	// First message should be system init.
	if _, ok := messages[0].(*SystemMessage); !ok {
		t.Errorf("first message is %T, want *SystemMessage", messages[0])
	}

	// Should have captured the session ID.
	if session.SessionID() != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("SessionID = %q", session.SessionID())
	}
}

func TestSessionSendEmptyPrompt(t *testing.T) {
	s := NewSession()
	defer func() { _ = s.Close() }()

	err := s.Send(context.Background(), "")
	if err != ErrEmptyPrompt {
		t.Errorf("expected ErrEmptyPrompt, got %v", err)
	}
}

func TestSessionSendAfterClose(t *testing.T) {
	s := NewSession()
	_ = s.Close()

	err := s.Send(context.Background(), "test")
	if err != ErrSessionClosed {
		t.Errorf("expected ErrSessionClosed, got %v", err)
	}
}

func TestSessionStreamAfterClose(t *testing.T) {
	s := NewSession()
	_ = s.Close()

	for _, err := range s.Stream(context.Background()) {
		if err != ErrSessionClosed {
			t.Errorf("expected ErrSessionClosed, got %v", err)
		}
		return
	}
}

func TestSessionCLINotFound(t *testing.T) {
	s := NewSession(WithCLIPath("/nonexistent/claude"))
	defer func() { _ = s.Close() }()

	err := s.Send(context.Background(), "test")
	if err != ErrCLINotFound {
		t.Errorf("expected ErrCLINotFound, got %v", err)
	}
}

func TestSessionDoubleClose(t *testing.T) {
	s := NewSession()
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Errorf("double close should not error, got: %v", err)
	}
}
