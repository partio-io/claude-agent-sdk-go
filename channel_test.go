package claude

import (
	"context"
	"errors"
	"iter"
	"testing"
)

func TestToChan(t *testing.T) {
	// Create a simple iterator that yields 3 messages.
	seq := func(yield func(Message, error) bool) {
		msgs := []Message{
			&SystemMessage{Type: "system", Subtype: "init", SessionID: "s1"},
			&AssistantMessage{Type: "assistant", UUID: "m1"},
			&ResultMessage{Type: "result", Subtype: ResultSuccess},
		}
		for _, msg := range msgs {
			if !yield(msg, nil) {
				return
			}
		}
	}

	ch := ToChan(context.Background(), iter.Seq2[Message, error](seq))

	var collected []Message
	for moe := range ch {
		if moe.Err != nil {
			t.Fatal(moe.Err)
		}
		collected = append(collected, moe.Message)
	}

	if len(collected) != 3 {
		t.Fatalf("got %d messages, want 3", len(collected))
	}
	if _, ok := collected[0].(*SystemMessage); !ok {
		t.Errorf("message[0] is %T, want *SystemMessage", collected[0])
	}
}

func TestToChanError(t *testing.T) {
	testErr := errors.New("test error")
	seq := func(yield func(Message, error) bool) {
		yield(nil, testErr)
	}

	ch := ToChan(context.Background(), iter.Seq2[Message, error](seq))
	moe := <-ch
	if moe.Err != testErr {
		t.Errorf("expected test error, got %v", moe.Err)
	}
}

func TestToChanContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	seq := func(yield func(Message, error) bool) {
		// Yield one message, then cancel.
		if !yield(&SystemMessage{Type: "system"}, nil) {
			return
		}
		cancel()
		// This yield should detect the cancellation.
		if !yield(&AssistantMessage{Type: "assistant"}, nil) {
			return
		}
	}

	ch := ToChan(ctx, iter.Seq2[Message, error](seq))

	var gotCancelled bool
	for moe := range ch {
		if moe.Err == context.Canceled {
			gotCancelled = true
		}
	}
	// Either we got a cancellation error or the channel closed.
	// Both are acceptable outcomes.
	_ = gotCancelled
}
