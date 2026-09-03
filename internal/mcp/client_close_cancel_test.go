package mcp

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestClientCloseContextCancelsConnectionInFlight(t *testing.T) {
	started := make(chan struct{})
	callDone := make(chan error, 1)

	client := NewClient(func(ctx context.Context) (Session, error) {
		close(started)
		<-ctx.Done()
		return nil, ctx.Err()
	})

	go func() {
		_, err := client.Tools(context.Background())
		callDone <- err
	}()

	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := client.CloseContext(ctx); err != nil {
		t.Fatalf("CloseContext: %v", err)
	}

	select {
	case err := <-callDone:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Tools error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Tools did not return after CloseContext")
	}
}
