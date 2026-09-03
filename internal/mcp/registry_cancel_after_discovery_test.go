package mcp

import (
	"context"
	"errors"
	"testing"
)

func TestRegistryDiscardsDiscoveryReturnedAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		calls++
		cancel()
		return []Tool{{Name: "read_script"}}, nil
	})

	tools, err := registry.Tools(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Tools error = %v, want context.Canceled", err)
	}
	if tools != nil {
		t.Fatalf("Tools = %v, want nil after cancellation", tools)
	}

	fresh, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("fresh Tools: %v", err)
	}
	if calls != 2 {
		t.Fatalf("discovery calls = %d, want 2 because canceled result must not be cached", calls)
	}
	if len(fresh) != 1 || fresh[0].Name != "read_script" {
		t.Fatalf("fresh tools = %v", fresh)
	}
}
