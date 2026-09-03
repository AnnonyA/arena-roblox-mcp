package mcp

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestRegistryWaiterRetriesAfterOwnerCancellation(t *testing.T) {
	started := make(chan struct{})
	var mu sync.Mutex
	calls := 0
	registry := NewRegistry(func(ctx context.Context) ([]Tool, error) {
		mu.Lock()
		calls++
		call := calls
		mu.Unlock()
		if call == 1 {
			close(started)
			<-ctx.Done()
			return nil, ctx.Err()
		}
		return []Tool{{Name: "read_script"}}, nil
	})

	ownerCtx, cancelOwner := context.WithCancel(context.Background())
	ownerDone := make(chan error, 1)
	go func() { _, err := registry.Tools(ownerCtx); ownerDone <- err }()
	<-started

	waiterDone := make(chan error, 1)
	go func() { _, err := registry.Tools(context.Background()); waiterDone <- err }()
	time.Sleep(10 * time.Millisecond)
	cancelOwner()

	if err := <-ownerDone; !errors.Is(err, context.Canceled) {
		t.Fatalf("owner Tools error = %v, want context.Canceled", err)
	}
	select {
	case err := <-waiterDone:
		if err != nil {
			t.Fatalf("waiter Tools error = %v, want nil", err)
		}
	case <-time.After(time.Second):
		t.Fatal("waiter did not retry discovery after owner cancellation")
	}
	mu.Lock()
	gotCalls := calls
	mu.Unlock()
	if gotCalls != 2 {
		t.Fatalf("discover calls = %d, want 2", gotCalls)
	}
}
