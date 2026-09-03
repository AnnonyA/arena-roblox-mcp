package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
)

type cancelAfterConnectSession struct {
	closes atomic.Int32
}

func (*cancelAfterConnectSession) ListTools(context.Context) ([]Tool, error) {
	return []Tool{{Name: "read_script"}}, nil
}

func (*cancelAfterConnectSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{}, nil
}

func (s *cancelAfterConnectSession) Close() error {
	s.closes.Add(1)
	return nil
}

func TestClientDiscardsSessionReturnedAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stale := &cancelAfterConnectSession{}
	fresh := &cancelAfterConnectSession{}
	connectCalls := 0

	client := NewClient(func(context.Context) (Session, error) {
		connectCalls++
		cancel()
		return stale, nil
	})

	_, err := client.Tools(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Tools error = %v, want context.Canceled", err)
	}
	if got := stale.closes.Load(); got != 1 {
		t.Fatalf("stale session Close calls = %d, want 1", got)
	}

	client.connect = func(context.Context) (Session, error) {
		connectCalls++
		return fresh, nil
	}
	if _, err := client.Tools(context.Background()); err != nil {
		t.Fatalf("Tools after cancellation: %v", err)
	}
	if connectCalls != 2 {
		t.Fatalf("connect calls = %d, want 2", connectCalls)
	}
}
