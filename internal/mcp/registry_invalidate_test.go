package mcp

import (
	"context"
	"testing"
)

func TestRegistryInvalidateForcesRediscovery(t *testing.T) {
	calls := 0
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		calls++
		return []Tool{{Name: "tool-v" + string(rune('0'+calls))}}, nil
	})

	first, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("first Tools: %v", err)
	}
	registry.Invalidate()
	second, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("second Tools: %v", err)
	}

	if calls != 2 {
		t.Fatalf("discovery calls = %d, want 2", calls)
	}
	if first[0].Name == second[0].Name {
		t.Fatalf("invalidate returned stale tools: first=%q second=%q", first[0].Name, second[0].Name)
	}
}

func TestRegistryInvalidateDuringDiscoveryDoesNotCacheStaleResult(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	firstDone := make(chan error, 1)
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
		_, err := registry.Tools(context.Background())
		firstDone <- err
	}()

	<-started
	registry.Invalidate()
	close(release)
	if err := <-firstDone; err != nil {
		t.Fatalf("first Tools: %v", err)
	}

	tools, err := registry.Tools(context.Background())
	if err != nil {
		t.Fatalf("second Tools: %v", err)
	}
	if calls != 2 {
		t.Fatalf("discovery calls = %d, want 2", calls)
	}
	if len(tools) != 1 || tools[0].Name != "fresh" {
		t.Fatalf("Tools after invalidate = %v, want fresh", tools)
	}
}
