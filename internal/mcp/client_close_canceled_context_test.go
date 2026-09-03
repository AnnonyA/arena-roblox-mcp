package mcp

import (
	"context"
	"errors"
	"testing"
)

func TestClientCloseContextWithCanceledContextStillClosesClient(t *testing.T) {
	client := NewClient(func(context.Context) (Session, error) {
		t.Fatal("connector should not be called after close")
		return nil, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := client.CloseContext(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("CloseContext error = %v, want context.Canceled", err)
	}

	if _, err := client.Tools(context.Background()); !errors.Is(err, ErrClientClosed) {
		t.Fatalf("Tools after CloseContext error = %v, want ErrClientClosed", err)
	}
}
