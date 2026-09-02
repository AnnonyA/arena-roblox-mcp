package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestRegistryCachesSuccessfulDiscovery(t *testing.T) {
	calls := 0
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		calls++
		return []Tool{{
			Name:        "read_script",
			Description: "Read Luau source",
			InputSchema: json.RawMessage(`{"type":"object"}`),
		}}, nil
	})

	first, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("first Tools: %v", err)
	}
	second, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("second Tools: %v", err)
	}

	if calls != 1 {
		t.Fatalf("discovery calls = %d, want 1", calls)
	}
	if len(first) != 1 || len(second) != 1 || first[0].Name != "read_script" || second[0].Name != "read_script" {
		t.Fatalf("unexpected tools: first=%v second=%v", first, second)
	}
}

func TestRegistryRetriesAfterDiscoveryFailure(t *testing.T) {
	calls := 0
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		calls++
		if calls == 1 {
			return nil, errors.New("studio unavailable")
		}
		return []Tool{{Name: "search_scripts"}}, nil
	})

	if _, err := registry.Tools(context.Background()); err == nil {
		t.Fatal("first Tools error = nil, want discovery failure")
	}
	tools, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("second Tools: %v", err)
	}
	if calls != 2 {
		t.Fatalf("discovery calls = %d, want 2", calls)
	}
	if len(tools) != 1 || tools[0].Name != "search_scripts" {
		t.Fatalf("unexpected tools: %v", tools)
	}
}

func TestRegistryReturnsDefensiveSchemaCopies(t *testing.T) {
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return []Tool{{Name: "read_script", InputSchema: json.RawMessage(`{"type":"object"}`)}}, nil
	})

	first, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("first Tools: %v", err)
	}
	first[0].InputSchema[0] = '['

	second, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("second Tools: %v", err)
	}
	if string(second[0].InputSchema) != `{"type":"object"}` {
		t.Fatalf("cached schema mutated: %q", second[0].InputSchema)
	}
}
