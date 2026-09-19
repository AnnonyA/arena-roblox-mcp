package cli

import (
	"strings"
	"testing"
)

func BenchmarkBoundedSafeDisplayTextShort(b *testing.B) {
	const input = "tool result: ready"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = boundedSafeDisplayText(input, maxToolDisplayRunes)
	}
}

func BenchmarkBoundedSafeMultilineDisplayTextShort(b *testing.B) {
	const input = "arena:\n  model: test-model\n"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = boundedSafeMultilineDisplayText(input, maxConfigDisplayRunes)
	}
}

func BenchmarkBoundedSafeMultilineDisplayTextTruncated(b *testing.B) {
	input := strings.Repeat("x", maxDiffDisplayRunes*2)

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		_ = boundedSafeMultilineDisplayText(input, maxDiffDisplayRunes)
	}
}

func BenchmarkBoundedSafeDisplayTextFilteredPrefix(b *testing.B) {
	// Exercise the security-filtering path before visible output reaches its
	// bound. This guards against regressions where filtered input causes
	// avoidable allocations or repeated rescans.
	input := strings.Repeat("\x1b\u202e", maxToolDisplayRunes) + strings.Repeat("x", maxToolDisplayRunes*2)

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		_ = boundedSafeDisplayText(input, maxToolDisplayRunes)
	}
}
