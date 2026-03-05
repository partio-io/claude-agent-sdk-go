package claude

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalContentBlock(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantType string
	}{
		{
			name:     "text",
			json:     `{"type":"text","text":"Hello world"}`,
			wantType: "text",
		},
		{
			name:     "thinking",
			json:     `{"type":"thinking","thinking":"Let me think...","signature":"sig_abc"}`,
			wantType: "thinking",
		},
		{
			name:     "tool_use",
			json:     `{"type":"tool_use","id":"toolu_01","name":"Read","input":{"file_path":"/tmp/test"}}`,
			wantType: "tool_use",
		},
		{
			name:     "tool_result",
			json:     `{"type":"tool_result","tool_use_id":"toolu_01","content":"file contents","is_error":false}`,
			wantType: "tool_result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := unmarshalContentBlock(json.RawMessage(tt.json))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := block.blockType(); got != tt.wantType {
				t.Errorf("blockType() = %q, want %q", got, tt.wantType)
			}
		})
	}
}

func TestUnmarshalTextBlock(t *testing.T) {
	raw := json.RawMessage(`{"type":"text","text":"Hello world"}`)
	block, err := unmarshalContentBlock(raw)
	if err != nil {
		t.Fatal(err)
	}
	tb, ok := block.(*TextBlock)
	if !ok {
		t.Fatalf("expected *TextBlock, got %T", block)
	}
	if tb.Text != "Hello world" {
		t.Errorf("Text = %q, want %q", tb.Text, "Hello world")
	}
}

func TestUnmarshalToolUseBlock(t *testing.T) {
	raw := json.RawMessage(`{"type":"tool_use","id":"toolu_01","name":"Read","input":{"file_path":"/tmp/test"}}`)
	block, err := unmarshalContentBlock(raw)
	if err != nil {
		t.Fatal(err)
	}
	tub, ok := block.(*ToolUseBlock)
	if !ok {
		t.Fatalf("expected *ToolUseBlock, got %T", block)
	}
	if tub.Name != "Read" {
		t.Errorf("Name = %q, want %q", tub.Name, "Read")
	}
	if tub.ID != "toolu_01" {
		t.Errorf("ID = %q, want %q", tub.ID, "toolu_01")
	}
	if tub.Input["file_path"] != "/tmp/test" {
		t.Errorf("Input[file_path] = %v, want /tmp/test", tub.Input["file_path"])
	}
}

func TestUnmarshalThinkingBlock(t *testing.T) {
	raw := json.RawMessage(`{"type":"thinking","thinking":"Let me think...","signature":"sig_abc"}`)
	block, err := unmarshalContentBlock(raw)
	if err != nil {
		t.Fatal(err)
	}
	tb, ok := block.(*ThinkingBlock)
	if !ok {
		t.Fatalf("expected *ThinkingBlock, got %T", block)
	}
	if tb.Thinking != "Let me think..." {
		t.Errorf("Thinking = %q, want %q", tb.Thinking, "Let me think...")
	}
	if tb.Signature != "sig_abc" {
		t.Errorf("Signature = %q, want %q", tb.Signature, "sig_abc")
	}
}

func TestUnmarshalContentBlocks(t *testing.T) {
	raw := json.RawMessage(`[{"type":"text","text":"Hello"},{"type":"tool_use","id":"t1","name":"Bash","input":{}}]`)
	blocks, err := unmarshalContentBlocks(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2", len(blocks))
	}
	if blocks[0].blockType() != "text" {
		t.Errorf("block[0] type = %q, want text", blocks[0].blockType())
	}
	if blocks[1].blockType() != "tool_use" {
		t.Errorf("block[1] type = %q, want tool_use", blocks[1].blockType())
	}
}

func TestUnmarshalContentBlockUnknownType(t *testing.T) {
	raw := json.RawMessage(`{"type":"unknown_type","data":"value"}`)
	block, err := unmarshalContentBlock(raw)
	if err != nil {
		t.Fatal(err)
	}
	// Unknown types fall back to TextBlock.
	tb, ok := block.(*TextBlock)
	if !ok {
		t.Fatalf("expected *TextBlock for unknown type, got %T", block)
	}
	if tb.Type != "unknown_type" {
		t.Errorf("Type = %q, want unknown_type", tb.Type)
	}
}

func TestContentBlockSealedInterface(t *testing.T) {
	// Verify all types implement ContentBlock.
	var _ ContentBlock = (*TextBlock)(nil)
	var _ ContentBlock = (*ThinkingBlock)(nil)
	var _ ContentBlock = (*ToolUseBlock)(nil)
	var _ ContentBlock = (*ToolResultBlock)(nil)
}
