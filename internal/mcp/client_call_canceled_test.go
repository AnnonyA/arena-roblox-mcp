package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type canceledCallSession struct{ calls int }

func (s *canceledCallSession) ListTools(context.Context) ([]Tool, error) { return nil, nil }
func (s *canceledCallSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	s.calls++
	return ToolResult{}, nil
}
func (s *canceledCallSession) Close() error { return nil }

func TestClientCallToolDoesNotInvokeSessionAfterCancellation(t *testing.T) {
	session := &canceledCallSession{}
	client := NewClient(func(context.Context) (Session, error) { return session, nil })

	if _, err := client.CallTool(context.Background(), "read_script", nil); err != nil {
		t.Fatalf("initial CallTool: %v", err)
	}
	session.calls = 0

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := client.CallTool(ctx, "read_script", nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("CallTool error = %v, want context.Canceled", err)
	}
	if session.calls != 0 {
		t.Fatalf("session CallTool calls = %d, want 0", session.calls)
	}
}
