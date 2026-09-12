package mcp

import "testing"

func TestValidateToolNamesRejectsUnicodeFormatCharacters(t *testing.T) {
	tests := []string{
		"read\u200bscript",
		"read\u202escript",
		"read\u2066script",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			if err := validateToolNames([]Tool{{Name: name}}); err == nil {
				t.Fatalf("validateToolNames(%q) error = nil, want Unicode format-character rejection", name)
			}
		})
	}
}
