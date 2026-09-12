package mcp

import (
	"context"
	"strings"
	"testing"
)

func TestRegistryRejectsBlankToolNames(t *testing.T) {
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return []Tool{
			{Name: "read_script"},
			{Name: " \t "},
		}, nil
	})

	tools, err := registry.Tools(context.Background())
	if err == nil {
		t.Fatal("Tools() error = nil, want blank tool name rejection")
	}
	if tools != nil {
		t.Fatalf("Tools() = %v, want nil on blank tool name", tools)
	}
	if !strings.Contains(err.Error(), "blank MCP tool name") {
		t.Fatalf("Tools() error = %q, want blank MCP tool name", err)
	}
}
