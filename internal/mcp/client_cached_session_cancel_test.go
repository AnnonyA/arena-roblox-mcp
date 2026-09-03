package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type cachedSession struct{}

func (cachedSession) ListTools(context.Context) ([]Tool, error) { return nil, nil }

func (cachedSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{}, nil
}

func (cachedSession) Close() error { return nil }

func TestEnsureSessionRejectsCanceledContextWithCachedSession(t *testing.T) {
	client := NewClient(nil)
	client.session = cachedSession{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.ensureSession(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ensureSession error = %v, want context.Canceled", err)
	}
}
