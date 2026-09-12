package mcp

import (
	"context"
	"strings"
	"testing"
)

func TestRegistryRejectsToolNamesWithSurroundingWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
	}{
		{name: "leading space", toolName: " read_script"},
		{name: "trailing space", toolName: "read_script "},
		{name: "leading tab", toolName: "\tread_script"},
		{name: "trailing newline", toolName: "read_script\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewRegistry(func(context.Context) ([]Tool, error) {
				return []Tool{{Name: tt.toolName}}, nil
			})

			_, err := registry.Tools(context.Background())
			if err == nil {
				t.Fatalf("Tools() error = nil for %q, want surrounding whitespace rejection", tt.toolName)
			}
			if !strings.Contains(err.Error(), "surrounding whitespace") {
				t.Fatalf("Tools() error = %q, want surrounding whitespace error", err)
			}
		})
	}
}
