package claude

import (
	"testing"

	"github.com/anthropics/claude-agent-sdk-go/internal/testutil"
)

func BenchmarkUnmarshalMessage(b *testing.B) {
	benchmarks := []struct {
		name string
		data string
	}{
		{"System", testutil.SystemInit},
		{"AssistantText", testutil.AssistantTextOnly},
		{"AssistantToolUse", testutil.AssistantToolUse},
		{"AssistantThinking", testutil.AssistantThinking},
		{"User", testutil.UserToolResult},
		{"ResultSuccess", testutil.ResultSuccess},
		{"StreamEvent", testutil.StreamEventTextDelta},
	}

	for _, bm := range benchmarks {
		data := []byte(bm.data)
		b.Run(bm.name, func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			b.ReportAllocs()
			for b.Loop() {
				_, err := unmarshalMessage(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
