package protocol

import "testing"

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string
	}{
		{"system", `{"type":"system","subtype":"init"}`, "system"},
		{"assistant", `{"type":"assistant","uuid":"msg1"}`, "assistant"},
		{"user", `{"type":"user","uuid":"usr1"}`, "user"},
		{"result", `{"type":"result","subtype":"success"}`, "result"},
		{"stream_event", `{"type":"stream_event","event":{}}`, "stream_event"},
		{"control_request", `{"type":"control_request","request_id":"r1"}`, "control_request"},
		{"control_response", `{"type":"control_response","response":{}}`, "control_response"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, err := ParseLine([]byte(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			if rm.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", rm.Type, tt.wantType)
			}
		})
	}
}

func TestParseLineInvalid(t *testing.T) {
	_, err := ParseLine([]byte(`{invalid`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestIsControlRequest(t *testing.T) {
	rm := &RawMessage{Type: "control_request"}
	if !IsControlRequest(rm) {
		t.Error("should be control request")
	}
	rm.Type = "assistant"
	if IsControlRequest(rm) {
		t.Error("should not be control request")
	}
}

func TestIsMessage(t *testing.T) {
	messageTypes := []string{"system", "assistant", "user", "result", "stream_event"}
	for _, mt := range messageTypes {
		rm := &RawMessage{Type: mt}
		if !IsMessage(rm) {
			t.Errorf("%q should be a message type", mt)
		}
	}

	nonMessageTypes := []string{"control_request", "control_response", "unknown"}
	for _, mt := range nonMessageTypes {
		rm := &RawMessage{Type: mt}
		if IsMessage(rm) {
			t.Errorf("%q should not be a message type", mt)
		}
	}
}

func TestParseLineCopiesData(t *testing.T) {
	buf := []byte(`{"type":"system"}`)
	rm, err := ParseLine(buf)
	if err != nil {
		t.Fatal(err)
	}
	// Mutate original buffer.
	buf[0] = 'X'
	// RawMessage data should be unaffected.
	if rm.Data[0] != '{' {
		t.Error("ParseLine should copy data")
	}
}
