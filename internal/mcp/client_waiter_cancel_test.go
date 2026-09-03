package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type retrySession struct{}

func (retrySession) ListTools(context.Context) ([]Tool, error) { return nil, nil }
func (retrySession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{Content: json.RawMessage(`{"ok":true}`)}, nil
}
func (retrySession) Close() error { return nil }

func TestClientWaiterRetriesWhenConnectionOwnerCancels(t *testing.T) {
	started := make(chan struct{})
	var connectCalls atomic.Int32

	client := NewClient(func(ctx context.Context) (Session, error) {
		switch connectCalls.Add(1) {
		case 1:
			close(started)
			<-ctx.Done()
			return nil, ctx.Err()
		case 2:
			return retrySession{}, nil
		default:
			t.Fatalf("unexpected extra connect call")
			return nil, errors.New("unreachable")
		}
	})

	ownerCtx, cancelOwner := context.WithCancel(context.Background())
	ownerErr := make(chan error, 1)
	go func() {
		_, err := client.CallTool(ownerCtx, "read_script", nil)
		ownerErr <- err
	}()
	<-started

	waiterErr := make(chan error, 1)
	go func() {
		_, err := client.CallTool(context.Background(), "read_script", nil)
		waiterErr <- err
	}()

	time.Sleep(20 * time.Millisecond)
	cancelOwner()

	if err := <-ownerErr; !errors.Is(err, context.Canceled) {
		t.Fatalf("owner error = %v, want context.Canceled", err)
	}

	select {
	case err := <-waiterErr:
		if err != nil {
			t.Fatalf("waiter error = %v, want nil", err)
		}
	case <-time.After(time.Second):
		t.Fatal("waiter did not recover from owner cancellation")
	}

	if got := connectCalls.Load(); got != 2 {
		t.Fatalf("connect calls = %d, want 2", got)
	}
}
