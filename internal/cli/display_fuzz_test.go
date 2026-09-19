package cli

import (
	"testing"
	"unicode"
	"unicode/utf8"
)

func FuzzBoundedSafeDisplayText(f *testing.F) {
	f.Add("plain text", uint16(32))
	f.Add("line one\nline two\tvalue", uint16(64))
	f.Add("\x1b[31mred\x1b[0m\u202Ehidden", uint16(16))
	f.Add("こんにちは世界", uint16(5))

	f.Fuzz(func(t *testing.T, input string, rawLimit uint16) {
		maxRunes := int(rawLimit % 4097)
		got := boundedSafeDisplayText(input, maxRunes)
		if count := utf8.RuneCountInString(got); count > maxRunes {
			t.Fatalf("output has %d runes, limit is %d", count, maxRunes)
		}
		for _, r := range got {
			if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
				t.Fatalf("output contains unsafe rune U+%04X", r)
			}
		}
	})
}
