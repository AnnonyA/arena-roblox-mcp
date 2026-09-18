package cli

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func FuzzBoundedSafeMultilineDisplayText(f *testing.F) {
	f.Add("plain text", 32)
	f.Add("line one\nline two\tvalue", 64)
	f.Add("\x1b[31mred\x1b[0m\u202Ehidden", 16)
	f.Add(strings.Repeat("界", 128), 12)

	f.Fuzz(func(t *testing.T, input string, maxRunes int) {
		if maxRunes < 0 || maxRunes > 4096 {
			return
		}
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
