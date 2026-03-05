package claude

import (
	"testing"

	"github.com/anthropics/claude-agent-sdk-go/internal/testutil"
)

func TestUnmarshalSystemMessage(t *testing.T) {
	msg, err := unmarshalMessage([]byte(testutil.SystemInit))
	if err != nil {
		t.Fatal(err)
	}
	sys, ok := msg.(*SystemMessage)
	if !ok {
		t.Fatalf("expected *SystemMessage, got %T", msg)
	}
	if sys.Subtype != "init" {
		t.Errorf("Subtype = %q, want init", sys.Subtype)
	}
	if sys.SessionID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("SessionID = %q", sys.SessionID)
	}
	if sys.Model != "claude-sonnet-4-6" {
		t.Errorf("Model = %q", sys.Model)
	}
	if len(sys.Tools) == 0 {
		t.Error("Tools should not be empty")
	}
}

func TestUnmarshalAssistantTextMessage(t *testing.T) {
	msg, err := unmarshalMessage([]byte(testutil.AssistantTextOnly))
	if err != nil {
		t.Fatal(err)
	}
	am, ok := msg.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", msg)
	}
	if am.UUID != "msg_01XFG" {
		t.Errorf("UUID = %q", am.UUID)
	}
	if am.Message == nil {
		t.Fatal("Message is nil")
	}
	if len(am.Message.Content) != 1 {
		t.Fatalf("Content length = %d, want 1", len(am.Message.Content))
	}
	tb, ok := am.Message.Content[0].(*TextBlock)
	if !ok {
		t.Fatalf("expected *TextBlock, got %T", am.Message.Content[0])
	}
	if tb.Text != "The answer is 4." {
		t.Errorf("Text = %q", tb.Text)
	}
	if am.Message.StopReason == nil || *am.Message.StopReason != "end_turn" {
		t.Errorf("StopReason = %v", am.Message.StopReason)
	}
}

func TestUnmarshalAssistantToolUseMessage(t *testing.T) {
	msg, err := unmarshalMessage([]byte(testutil.AssistantToolUse))
	if err != nil {
		t.Fatal(err)
	}
	am, ok := msg.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", msg)
	}
	if len(am.Message.Content) != 2 {
		t.Fatalf("Content length = %d, want 2", len(am.Message.Content))
	}

	// First block: text
	if _, ok := am.Message.Content[0].(*TextBlock); !ok {
		t.Errorf("Content[0] is %T, want *TextBlock", am.Message.Content[0])
	}

	// Second block: tool_use
	tub, ok := am.Message.Content[1].(*ToolUseBlock)
	if !ok {
		t.Fatalf("Content[1] is %T, want *ToolUseBlock", am.Message.Content[1])
	}
	if tub.Name != "Read" {
		t.Errorf("Name = %q, want Read", tub.Name)
	}
}

func TestUnmarshalAssistantThinkingMessage(t *testing.T) {
	msg, err := unmarshalMessage([]byte(testutil.AssistantThinking))
	if err != nil {
		t.Fatal(err)
	}
	am, ok := msg.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", msg)
	}
	if len(am.Message.Content) != 2 {
		t.Fatalf("Content length = %d, want 2", len(am.Message.Content))
	}
	tb, ok := am.Message.Content[0].(*ThinkingBlock)
	if !ok {
		t.Fatalf("Content[0] is %T, want *ThinkingBlock", am.Message.Content[0])
	}
	if tb.Thinking != "Let me analyze this..." {
		t.Errorf("Thinking = %q", tb.Thinking)
	}
}

func TestUnmarshalUserMessage(t *testing.T) {
	msg, err := unmarshalMessage([]byte(testutil.UserToolResult))
	if err != nil {
		t.Fatal(err)
	}
	um, ok := msg.(*UserMessage)
	if !ok {
		t.Fatalf("expected *UserMessage, got %T", msg)
	}
	if um.UUID != "usr_01DEF" {
		t.Errorf("UUID = %q", um.UUID)
	}
	if um.Message == nil {
		t.Fatal("Message is nil")
	}
	if len(um.Message.Content) != 1 {
		t.Fatalf("Content length = %d, want 1", len(um.Message.Content))
	}
	tr, ok := um.Message.Content[0].(*ToolResultBlock)
	if !ok {
		t.Fatalf("expected *ToolResultBlock, got %T", um.Message.Content[0])
	}
	if tr.ToolUseID != "toolu_01ABC" {
		t.Errorf("ToolUseID = %q", tr.ToolUseID)
	}
}

func TestUnmarshalResultSuccess(t *testing.T) {
	msg, err := unmarshalMessage([]byte(testutil.ResultSuccess))
	if err != nil {
		t.Fatal(err)
	}
	rm, ok := msg.(*ResultMessage)
	if !ok {
		t.Fatalf("expected *ResultMessage, got %T", msg)
	}
	if rm.Subtype != ResultSuccess {
		t.Errorf("Subtype = %q, want success", rm.Subtype)
	}
	if rm.IsError {
		t.Error("IsError should be false")
	}
	if rm.Result == nil || *rm.Result != "The answer is 4." {
		t.Errorf("Result = %v", rm.Result)
	}
	if rm.NumTurns != 1 {
		t.Errorf("NumTurns = %d, want 1", rm.NumTurns)
	}
	if rm.TotalCostUSD == nil || *rm.TotalCostUSD != 0.0023 {
		t.Errorf("TotalCostUSD = %v", rm.TotalCostUSD)
	}
}

func TestUnmarshalResultMaxTurns(t *testing.T) {
	msg, err := unmarshalMessage([]byte(testutil.ResultMaxTurns))
	if err != nil {
		t.Fatal(err)
	}
	rm, ok := msg.(*ResultMessage)
	if !ok {
		t.Fatalf("expected *ResultMessage, got %T", msg)
	}
	if rm.Subtype != ResultErrorMaxTurns {
		t.Errorf("Subtype = %q, want error_max_turns", rm.Subtype)
	}
	if !rm.IsError {
		t.Error("IsError should be true")
	}
}

func TestUnmarshalStreamEvent(t *testing.T) {
	msg, err := unmarshalMessage([]byte(testutil.StreamEventTextDelta))
	if err != nil {
		t.Fatal(err)
	}
	se, ok := msg.(*StreamEvent)
	if !ok {
		t.Fatalf("expected *StreamEvent, got %T", msg)
	}
	if se.UUID != "evt_01GHI" {
		t.Errorf("UUID = %q", se.UUID)
	}
	if se.Event["type"] != "content_block_delta" {
		t.Errorf("Event type = %v", se.Event["type"])
	}
}

func TestUnmarshalUnknownType(t *testing.T) {
	_, err := unmarshalMessage([]byte(`{"type":"unknown"}`))
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	pe, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Raw == nil {
		t.Error("ParseError.Raw should not be nil")
	}
}

func TestUnmarshalInvalidJSON(t *testing.T) {
	_, err := unmarshalMessage([]byte(`{invalid json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestMessageSealedInterface(t *testing.T) {
	// Verify all types implement Message.
	var _ Message = (*SystemMessage)(nil)
	var _ Message = (*AssistantMessage)(nil)
	var _ Message = (*UserMessage)(nil)
	var _ Message = (*ResultMessage)(nil)
	var _ Message = (*StreamEvent)(nil)
}
