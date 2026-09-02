package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"
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

func TestRegistryWaiterCanCancelDuringDiscovery(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	firstDone := make(chan error, 1)
	var calls atomic.Int32

	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		<-release
		return []Tool{{Name: "read_script"}}, nil
	})

	go func() {
		_, err := registry.Tools(context.Background())
		firstDone <- err
	}()

	<-started

	ctx, cancel := context.WithCancel(context.Background())
	secondDone := make(chan error, 1)
	go func() {
		_, err := registry.Tools(ctx)
		secondDone <- err
	}()
	cancel()

	select {
	case err := <-secondDone:
		if !errors.Is(err, context.Canceled) {
			close(release)
			<-firstDone
			t.Fatalf("second Tools error = %v, want context.Canceled", err)
		}
	case <-time.After(100 * time.Millisecond):
		close(release)
		<-firstDone
		t.Fatal("second Tools ignored cancellation while discovery was in progress")
	}

	close(release)
	if err := <-firstDone; err != nil {
		t.Fatalf("first Tools: %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("discovery calls = %d, want 1", got)
	}
}
