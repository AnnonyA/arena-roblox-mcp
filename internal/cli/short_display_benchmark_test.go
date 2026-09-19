package cli

import "testing"

func BenchmarkBoundedSafeMultilineDisplayTextShort(b *testing.B) {
	// Keep a representative short /diff payload covered. The multiline display
	// bound is intentionally large, so this benchmark makes oversized temporary
	// allocations on ordinary small outputs visible in benchmark reports.
	const input = "--- before\n+++ after\n-old\n+new\n"

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		_ = boundedSafeMultilineDisplayText(input, maxDiffDisplayRunes)
	}
}
