package cli

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func FuzzBoundedSafeMultilineDisplayText(f *testing.F) {
	f.Add("plain text", uint16(32))
	f.Add("line one\nline two\tvalue", uint16(64))
	f.Add("\x1b[31mred\x1b[0m\u202Ehidden", uint16(16))
	f.Add(strings.Repeat("界", 128), uint16(12))

	f.Fuzz(func(t *testing.T, input string, rawLimit uint16) {
		maxRunes := int(rawLimit % 4097)
		got := boundedSafeMultilineDisplayText(input, maxRunes)
		if utf8.RuneCountInString(got) > maxRunes {
			t.Fatalf("output has %d runes, limit is %d", utf8.RuneCountInString(got), maxRunes)
		}
		for _, r := range got {
			if r == '\n' || r == '\t' {
				continue
			}
			if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
				t.Fatalf("output contains unsafe rune %U", r)
			}
		}
	})
}
