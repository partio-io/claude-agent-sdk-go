package claude

import (
	"encoding/json"
	"testing"
)

func BenchmarkUnmarshalContentBlock(b *testing.B) {
	benchmarks := []struct {
		name string
		data string
	}{
		{"Text", `{"type":"text","text":"The answer is 4."}`},
		{"Thinking", `{"type":"thinking","thinking":"Let me analyze this...","signature":"sig_abc"}`},
		{"ToolUse", `{"type":"tool_use","id":"toolu_01ABC","name":"Read","input":{"file_path":"/home/user/project/main.go"}}`},
		{"ToolResult", `{"type":"tool_result","tool_use_id":"toolu_01ABC","content":"package main\n\nfunc main() {}","is_error":false}`},
	}

	for _, bm := range benchmarks {
		raw := json.RawMessage(bm.data)
		b.Run(bm.name, func(b *testing.B) {
			b.SetBytes(int64(len(raw)))
			b.ReportAllocs()
			for b.Loop() {
				_, err := unmarshalContentBlock(raw)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
