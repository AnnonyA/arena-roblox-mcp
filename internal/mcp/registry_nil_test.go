package mcp

import (
	"context"
	"errors"
	"testing"
)

func TestRegistryWithoutDiscovererReturnsConfigurationError(t *testing.T) {
	registry := NewRegistry(nil)

	_, err := registry.Tools(context.Background())
	if !errors.Is(err, ErrNoDiscoverer) {
		t.Fatalf("Tools error = %v, want ErrNoDiscoverer", err)
	}
}
