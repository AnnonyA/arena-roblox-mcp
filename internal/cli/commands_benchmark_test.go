package cli

import (
	"strings"
	"testing"
)

func BenchmarkUnknownCommandSuggestionLongInput(b *testing.B) {
	command := strings.Repeat("界", 64*1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got := suggestedCommand(command); got != "" {
			b.Fatalf("suggestedCommand() = %q, want empty", got)
		}
	}
}
