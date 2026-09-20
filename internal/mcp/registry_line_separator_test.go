package mcp

import (
	"context"
	"strings"
	"testing"
)

func TestRegistryRejectsLineSeparatorsInsideToolNames(t *testing.T) {
	for _, separator := range []string{"\u2028", "\u2029"} {
		t.Run(separator, func(t *testing.T) {
			registry := NewRegistry(func(context.Context) ([]Tool, error) {
				return []Tool{{Name: "read" + separator + "script"}}, nil
			})

			_, err := registry.Tools(context.Background())
			if err == nil {
				t.Fatal("expected unsafe MCP tool name to be rejected")
			}
			if !strings.Contains(err.Error(), "line separator") {
				t.Fatalf("expected line separator error, got %v", err)
			}
		})
	}
}
