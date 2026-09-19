package cli

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
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

func FuzzBoundedSafeMultilineDisplayTextIntLimit(f *testing.F) {
	for _, seed := range []struct {
		input string
		limit int
	}{
		{"line one\nline two\tvalue", 8},
		{"escape\x1b[2Jscreen", 6},
		{"hidden\u202etext", 4},
		{"🎮🎮🎮", 2},
		{"anything", 0},
	} {
		f.Add(seed.input, seed.limit)
	}

	f.Fuzz(func(t *testing.T, input string, limit int) {
		if limit < 0 || limit > 1024 {
			t.Skip()
		}
		got := boundedSafeMultilineDisplayText(input, limit)
		if utf8.RuneCountInString(got) > limit {
			t.Fatalf("output exceeds limit %d: %q", limit, got)
		}
		for _, r := range got {
			if r == '\n' || r == '\t' {
				continue
			}
			if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
				t.Fatalf("bounded output retained unsafe rune %U in %q", r, got)
			}
		}
	})
}
