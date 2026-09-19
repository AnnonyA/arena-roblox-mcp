package cli

import "testing"

func BenchmarkBoundedSafeMultilineDisplayTextShortDiff(b *testing.B) {
	// Keep a representative short /diff payload covered. The diff display
	// bound is intentionally much larger than ordinary output, so this benchmark
	// makes oversized temporary allocations on small diffs visible in reports.
	const input = "--- before\n+++ after\n-old\n+new\n"

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		_ = boundedSafeMultilineDisplayText(input, maxDiffDisplayRunes)
	}
}
