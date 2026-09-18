package cli

import (
	"testing"
	"unicode"
	"unicode/utf8"
)

func FuzzSafeModelIDs(f *testing.F) {
	for _, seed := range []string{
		"arena-code",
		"arena\x1b[31m-red",
		"arena\u202Ehidden",
		"模型-1",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, model string) {
		got := safeModelIDs([]string{model})

		wantSafe := utf8.RuneCountInString(model) <= maxModelIDDisplayRunes
		for _, r := range model {
			if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
				wantSafe = false
				break
			}
		}

		if !wantSafe {
			if len(got) != 0 {
				t.Fatalf("unsafe model ID was retained: %q", got)
			}
			return
		}
		if len(got) != 1 || got[0] != model {
			t.Fatalf("safe model ID changed: got %q, want %q", got, model)
		}
	})
}
