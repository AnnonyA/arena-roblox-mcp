package cli

import (
	"strings"
	"testing"
)

func BenchmarkBoundedSafeDisplayTextShort(b *testing.B) {
	const input = "tool result: ready"

	b.ReportAllocs()
	for b.Loop() {
		_ = boundedSafeDisplayText(input, maxToolDisplayRunes)
	}
}

func BenchmarkBoundedSafeMultilineDisplayTextShort(b *testing.B) {
	const input = "arena:\n  model: test-model\n"

	b.ReportAllocs()
	for b.Loop() {
		_ = boundedSafeMultilineDisplayText(input, maxConfigDisplayRunes)
	}
}

func BenchmarkBoundedSafeMultilineDisplayTextTruncated(b *testing.B) {
	input := strings.Repeat("x", maxDiffDisplayRunes*2)

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	for b.Loop() {
		_ = boundedSafeMultilineDisplayText(input, maxDiffDisplayRunes)
	}
}
