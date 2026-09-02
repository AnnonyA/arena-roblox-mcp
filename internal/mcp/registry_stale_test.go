package mcp

import (
	"context"
	"testing"
)

func TestRegistryInvalidateDuringDiscoveryRetriesCallerWithFreshTools(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	result := make(chan []Tool, 1)
	errs := make(chan error, 1)
	calls := 0

	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		calls++
		if calls == 1 {
			close(started)
			<-release
			return []Tool{{Name: "stale"}}, nil
		}
		return []Tool{{Name: "fresh"}}, nil
	})

	go func() {
		tools, err := registry.Tools(context.Background())
		if err != nil {
			errs <- err
			return
		}
		result <- tools
	}()

	<-started
	registry.Invalidate()
	close(release)

	select {
	case err := <-errs:
		t.Fatalf("Tools: %v", err)
	case tools := <-result:
		if len(tools) != 1 || tools[0].Name != "fresh" {
			t.Fatalf("Tools returned %v, want fresh tools after invalidate", tools)
		}
	}
	if calls != 2 {
		t.Fatalf("discovery calls = %d, want 2", calls)
	}
}
