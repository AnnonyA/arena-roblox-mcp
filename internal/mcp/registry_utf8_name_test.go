package mcp

import (
	"context"
	"strings"
	"testing"
)

func TestRegistryRejectsInvalidUTF8ToolName(t *testing.T) {
	invalidName := string([]byte{'r', 'e', 'a', 'd', '_', 0xff})
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return []Tool{{Name: invalidName}}, nil
	})

	_, err := registry.Tools(context.Background())
	if err == nil {
		t.Fatal("expected invalid UTF-8 tool name to be rejected")
	}
	if !strings.Contains(err.Error(), "invalid UTF-8") {
		t.Fatalf("expected invalid UTF-8 error, got %v", err)
	}
}
