package mcp

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCloseContextCancelsInFlightConnectEvenWhenShutdownContextCanceled(t *testing.T) {
	connectCanceled := make(chan struct{})
	started := make(chan struct{})
	client := NewClient(func(ctx context.Context) (Session, error) {
		close(started)
		<-ctx.Done()
		close(connectCanceled)
		return nil, ctx.Err()
	})

	connectDone := make(chan error, 1)
	go func() {
		_, err := client.Tools(context.Background())
		connectDone <- err
	}()
	<-started

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := client.CloseContext(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("CloseContext error = %v, want context.Canceled", err)
	}

	select {
	case <-connectCanceled:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("CloseContext did not cancel in-flight MCP connection")
	}

	if err := <-connectDone; !errors.Is(err, context.Canceled) {
		t.Fatalf("in-flight Tools error = %v, want context.Canceled", err)
	}
}
