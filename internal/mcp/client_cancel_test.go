package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestClientDoesNotConnectWhenContextAlreadyCanceled(t *testing.T) {
	connectCalls := 0
	client := NewClient(func(context.Context) (Session, error) {
		connectCalls++
		return &fakeCanceledSession{}, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Tools(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Tools error = %v, want context.Canceled", err)
	}
	if connectCalls != 0 {
		t.Fatalf("connect calls = %d, want 0", connectCalls)
	}
}

type fakeCanceledSession struct{}

func (*fakeCanceledSession) ListTools(context.Context) ([]Tool, error) { return nil, nil }
func (*fakeCanceledSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{}, nil
}
func (*fakeCanceledSession) Close() error { return nil }
