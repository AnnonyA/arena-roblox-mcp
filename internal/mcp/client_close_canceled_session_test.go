package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type trackedCloseSession struct{ closed bool }

func (s *trackedCloseSession) ListTools(context.Context) ([]Tool, error) {
	return []Tool{{Name: "read"}}, nil
}
func (s *trackedCloseSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{}, nil
}
func (s *trackedCloseSession) Close() error { s.closed = true; return nil }

func TestClientCloseContextCanceledStillClosesExistingSession(t *testing.T) {
	session := &trackedCloseSession{}
	client := NewClient(func(context.Context) (Session, error) { return session, nil })
	if _, err := client.Tools(context.Background()); err != nil {
		t.Fatalf("prime session: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := client.CloseContext(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("CloseContext error = %v, want context.Canceled", err)
	}
	if !session.closed {
		t.Fatal("existing session was not closed")
	}
}
