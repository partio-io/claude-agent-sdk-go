package claude

import "testing"

func BenchmarkNewUUID(b *testing.B) {
	b.SetBytes(16) // 16 bytes of entropy per UUID
	b.ReportAllocs()
	for b.Loop() {
		_ = newUUID()
	}
}
