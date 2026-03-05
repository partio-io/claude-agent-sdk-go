package protocol

import "testing"

// Realistic ~300-byte assistant message line for benchmarking.
var benchLine = []byte(`{"type":"assistant","uuid":"msg_01XFG","session_id":"550e8400-e29b-41d4-a716-446655440000","parent_tool_use_id":null,"message":{"id":"msg_01XFG","type":"message","role":"assistant","model":"claude-sonnet-4-6","content":[{"type":"text","text":"The answer is 4."}],"stop_reason":"end_turn","usage":{"input_tokens":100,"output_tokens":20}}}`)

func BenchmarkParseLine(b *testing.B) {
	b.SetBytes(int64(len(benchLine)))
	b.ReportAllocs()
	for b.Loop() {
		_, err := ParseLine(benchLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}
