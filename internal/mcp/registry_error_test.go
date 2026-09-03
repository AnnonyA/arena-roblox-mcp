package mcp

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestRegistrySharesDiscoveryFailureWithWaiters(t *testing.T) {
	wantErr := errors.New("studio unavailable")
	started := make(chan struct{})
	release := make(chan struct{})

	var mu sync.Mutex
	discoverCalls := 0
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		mu.Lock()
		discoverCalls++
		call := discoverCalls
		mu.Unlock()
		if call == 1 {
			close(started)
			<-release
		}
		return nil, wantErr
	})

	firstDone := make(chan error, 1)
	go func() {
		_, err := registry.Tools(context.Background())
		firstDone <- err
	}()

	<-started
	secondDone := make(chan error, 1)
	go func() {
		_, err := registry.Tools(context.Background())
		secondDone <- err
	}()

	select {
	case err := <-secondDone:
		t.Fatalf("second Tools returned before first discovery completed: %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	close(release)
	if err := <-firstDone; !errors.Is(err, wantErr) {
		t.Fatalf("first Tools error = %v, want %v", err, wantErr)
	}
	if err := <-secondDone; !errors.Is(err, wantErr) {
		t.Fatalf("second Tools error = %v, want %v", err, wantErr)
	}

	mu.Lock()
	gotCalls := discoverCalls
	mu.Unlock()
	if gotCalls != 1 {
		t.Fatalf("discover calls = %d, want 1 for concurrent waiters", gotCalls)
	}

	if _, err := registry.Tools(context.Background()); !errors.Is(err, wantErr) {
		t.Fatalf("retry Tools error = %v, want %v", err, wantErr)
	}
	mu.Lock()
	gotCalls = discoverCalls
	mu.Unlock()
	if gotCalls != 2 {
		t.Fatalf("discover calls after retry = %d, want 2", gotCalls)
	}
}
