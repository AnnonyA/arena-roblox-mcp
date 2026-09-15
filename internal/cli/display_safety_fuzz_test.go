package cli

import (
	"strings"
	"testing"
	"unicode"
)

func FuzzSafeDisplayText(f *testing.F) {
	for _, seed := range []string{
		"plain text",
		"escape\x1b[31mred",
		"hidden\u200btext",
		"line\nfeed\ttab",
		"emoji 🎮 and Luau",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		got := safeDisplayText(input)
		for _, r := range got {
			if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
				t.Fatalf("safeDisplayText(%q) retained unsafe rune %U in %q", input, r, got)
			}
		}
	})
}

func FuzzSafeMultilineDisplayText(f *testing.F) {
	for _, seed := range []string{
		"line one\nline two\tvalue",
		"escape\x1b[2Jscreen",
		"hidden\u202etext",
		"carriage\rreturn",
		"emoji 🎮 and Luau",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		got := safeMultilineDisplayText(input)
		if strings.Count(got, "\n") != strings.Count(input, "\n") {
			t.Fatalf("newline count changed: input %q output %q", input, got)
		}
		if strings.Count(got, "\t") != strings.Count(input, "\t") {
			t.Fatalf("tab count changed: input %q output %q", input, got)
		}
		for _, r := range got {
			if r == '\n' || r == '\t' {
				continue
			}
			if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
				t.Fatalf("safeMultilineDisplayText(%q) retained unsafe rune %U in %q", input, r, got)
			}
		}
	})
}
