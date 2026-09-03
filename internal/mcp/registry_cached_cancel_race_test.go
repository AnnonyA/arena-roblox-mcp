package mcp

import (
	"context"
	"errors"
	"testing"
	"time"
)

type observedCancelContext struct {
	context.Context
	checked chan struct{}
}

func (c *observedCancelContext) Err() error {
	err := c.Context.Err()
	select {
	case <-c.checked:
	default:
		close(c.checked)
	}
	return err
}

func TestRegistryToolsRejectsCancellationWhileWaitingForCacheLock(t *testing.T) {
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return []Tool{{Name: "read_script"}}, nil
	})
	if _, err := registry.Tools(context.Background()); err != nil {
		t.Fatalf("prime cache: %v", err)
	}

	base, cancel := context.WithCancel(context.Background())
	ctx := &observedCancelContext{Context: base, checked: make(chan struct{})}

	registry.mu.Lock()
	done := make(chan error, 1)
	go func() {
		_, err := registry.Tools(ctx)
		done <- err
	}()
	<-ctx.checked // Tools passed its pre-lock cancellation check.
	cancel()
	registry.mu.Unlock()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Tools error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Tools did not return")
	}
}
