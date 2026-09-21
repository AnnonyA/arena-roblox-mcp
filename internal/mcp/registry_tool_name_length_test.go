package mcp

import (
	"context"
	"strings"
	"testing"
)

func TestRegistryRejectsOversizedToolName(t *testing.T) {
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return []Tool{{Name: strings.Repeat("a", maxToolNameBytes+1)}}, nil
	})

	_, err := registry.Tools(context.Background())
	if err == nil {
		t.Fatal("expected oversized MCP tool name to be rejected")
	}
	if !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("expected size validation error, got %v", err)
	}
}

func TestRegistryAcceptsToolNameAtSizeLimit(t *testing.T) {
	name := strings.Repeat("a", maxToolNameBytes)
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return []Tool{{Name: name}}, nil
	})

	tools, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("expected tool name at size limit to be accepted: %v", err)
	}
	if len(tools) != 1 || tools[0].Name != name {
		t.Fatalf("unexpected tools: %#v", tools)
	}
}
