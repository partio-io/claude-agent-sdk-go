package protocol

import (
	"encoding/json"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
)

func TestNextRequestID(t *testing.T) {
	ctrl := NewController(nil)

	id1 := ctrl.NextRequestID()
	id2 := ctrl.NextRequestID()

	if !strings.HasPrefix(id1, "req_1_") {
		t.Errorf("first ID should start with req_1_, got %q", id1)
	}
	if !strings.HasPrefix(id2, "req_2_") {
		t.Errorf("second ID should start with req_2_, got %q", id2)
	}
	if id1 == id2 {
		t.Error("IDs should be unique")
	}
}

func TestControllerSendAndHandle(t *testing.T) {
	var written atomic.Value
	writeFn := func(v any) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		written.Store(b)
		return nil
	}

	ctrl := NewController(writeFn)

	// Start async send.
	done := make(chan *ControlResponseBody)
	go func() {
		resp, _ := ctrl.SendRequest("initialize", map[string]any{})
		done <- resp
	}()

	// Wait for the request to be written, then parse it.
	for written.Load() == nil {
		runtime.Gosched()
	}

	var req ControlRequest
	if err := json.Unmarshal(written.Load().([]byte), &req); err != nil {
		t.Fatal(err)
	}

	// Simulate response.
	ctrl.HandleResponse(ControlResponseBody{
		Subtype:   "success",
		RequestID: req.RequestID,
		Response:  json.RawMessage(`{"supported_commands":["initialize"]}`),
	})

	resp := <-done
	if resp.Subtype != "success" {
		t.Errorf("Subtype = %q, want success", resp.Subtype)
	}
}

func TestControllerSendResponse(t *testing.T) {
	var written []byte
	writeFn := func(v any) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		written = b
		return nil
	}

	ctrl := NewController(writeFn)
	err := ctrl.SendResponse("req_1_abc", "success", map[string]any{"behavior": "allow"})
	if err != nil {
		t.Fatal(err)
	}

	var resp ControlResponse
	if err := json.Unmarshal(written, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Response.Subtype != "success" {
		t.Errorf("Subtype = %q", resp.Response.Subtype)
	}
	if resp.Response.RequestID != "req_1_abc" {
		t.Errorf("RequestID = %q", resp.Response.RequestID)
	}
}

func TestControllerSendErrorResponse(t *testing.T) {
	var written []byte
	writeFn := func(v any) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		written = b
		return nil
	}

	ctrl := NewController(writeFn)
	err := ctrl.SendErrorResponse("req_2_def", "something went wrong")
	if err != nil {
		t.Fatal(err)
	}

	var resp ControlResponse
	if err := json.Unmarshal(written, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Response.Subtype != "error" {
		t.Errorf("Subtype = %q", resp.Response.Subtype)
	}
	if resp.Response.Error != "something went wrong" {
		t.Errorf("Error = %q", resp.Response.Error)
	}
}

func TestParseRequestSubtype(t *testing.T) {
	raw := json.RawMessage(`{"subtype":"can_use_tool","tool_name":"Bash"}`)
	st, err := ParseRequestSubtype(raw)
	if err != nil {
		t.Fatal(err)
	}
	if st != "can_use_tool" {
		t.Errorf("subtype = %q", st)
	}
}
