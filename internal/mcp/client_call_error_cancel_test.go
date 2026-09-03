package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type errorAfterCancelSession struct{ cancel context.CancelFunc }

func (s *errorAfterCancelSession) ListTools(context.Context) ([]Tool, error) { return nil, nil }
func (s *errorAfterCancelSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	s.cancel()
	return ToolResult{}, errors.New("transport closed")
}
func (s *errorAfterCancelSession) Close() error { return nil }

func TestClientCallToolPrefersCallerCancellationOverLateSessionError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	session := &errorAfterCancelSession{cancel: cancel}
	client := NewClient(func(context.Context) (Session, error) { return session, nil })

	_, err := client.CallTool(ctx, "read_script", nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("CallTool error = %v, want context.Canceled", err)
	}
}
