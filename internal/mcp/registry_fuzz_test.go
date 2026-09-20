package mcp

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func FuzzToolNameValidation(f *testing.F) {
	for _, seed := range []string{
		"inspect",
		" inspect",
		"inspect ",
		"ins\npect",
		"ins\u200bpect",
		"ins\u2060pect",
		"ins\u2028pect",
		"ins\u2029pect",
		"検査",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, name string) {
		if !utf8.ValidString(name) {
			t.Skip()
		}

		err := validateToolNames([]Tool{{Name: name}})
		unsafe := strings.TrimSpace(name) != name || strings.TrimSpace(name) == ""
		if !unsafe {
			for _, r := range name {
				if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || r == '\u2028' || r == '\u2029' {
					unsafe = true
					break
				}
			}
		}

		if unsafe && err == nil {
			t.Fatalf("validateToolNames accepted unsafe MCP tool name %q", name)
		}
		if !unsafe && err != nil {
			t.Fatalf("validateToolNames rejected valid MCP tool name %q: %v", name, err)
		}
	})
}
