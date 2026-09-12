package mcp

import (
	"context"
	"strings"
	"testing"
)

func TestRegistryRejectsDuplicateToolNames(t *testing.T) {
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return []Tool{
			{Name: "read_script"},
			{Name: "read_script"},
		}, nil
	})

	tools, err := registry.Tools(context.Background())
	if err == nil {
		t.Fatal("Tools() error = nil, want duplicate tool name rejection")
	}
	if tools != nil {
		t.Fatalf("Tools() = %v, want nil on duplicate tool name", tools)
	}
	if !strings.Contains(err.Error(), "duplicate MCP tool name") {
		t.Fatalf("Tools() error = %q, want duplicate MCP tool name", err)
	}
}
