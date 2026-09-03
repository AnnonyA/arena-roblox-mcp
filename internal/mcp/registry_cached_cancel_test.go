package mcp

import (
	"context"
	"errors"
	"testing"
)

func TestRegistryToolsHonorsCanceledContextWhenCached(t *testing.T) {
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return []Tool{{Name: "read_script"}}, nil
	})
	if _, err := registry.Tools(context.Background()); err != nil {
		t.Fatalf("initial Tools: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := registry.Tools(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("cached Tools error = %v, want context.Canceled", err)
	}
}
