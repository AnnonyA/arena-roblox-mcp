package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type partialErrorSession struct{ closed bool }

func (*partialErrorSession) ListTools(context.Context) ([]Tool, error) { return nil, nil }
func (*partialErrorSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{}, nil
}
func (s *partialErrorSession) Close() error { s.closed = true; return nil }

func TestClientClosesPartialSessionWhenConnectFails(t *testing.T) {
	wantErr := errors.New("connect failed")
	session := &partialErrorSession{}
	client := NewClient(func(context.Context) (Session, error) { return session, wantErr })

	_, err := client.CallTool(context.Background(), "read_script", nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("CallTool error = %v, want %v", err, wantErr)
	}
	if !session.closed {
		t.Fatal("partial session was not closed after connect error")
	}
}
