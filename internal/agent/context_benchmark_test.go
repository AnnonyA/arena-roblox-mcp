package agent

import (
	"strings"
	"testing"
)

func BenchmarkContextProcessing(b *testing.B) {
	payload := strings.Repeat("x", 4096)
	for i := 0; i < b.N; i++ {
		ctx := NewContext(100)
		for j := 0; j < 100; j++ {
			role := "assistant"
			content := "short message"
			if j%5 == 0 {
				role = "tool"
				content = payload
			}
			ctx.Add(Event{Role: role, Content: content})
		}
		ctx.CompactToolOutputs(512)
		_ = ctx.Events()
	}
}
