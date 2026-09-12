package mcp

import "testing"

func TestValidateToolNamesRejectsControlCharacters(t *testing.T) {
	tests := []string{
		"read\x00script",
		"read\x1bscript",
		"read\u007fscript",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			if err := validateToolNames([]Tool{{Name: name}}); err == nil {
				t.Fatalf("validateToolNames(%q) error = nil, want control-character rejection", name)
			}
		})
	}
}
