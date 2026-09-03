package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
)

type closeGuardSession struct {
	mu         sync.Mutex
	closeCalls int
}

func (s *closeGuardSession) ListTools(context.Context) ([]Tool, error) {
	return []Tool{{Name: "read_script"}}, nil
}

func (s *closeGuardSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{}, nil
}

func (s *closeGuardSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeCalls++
	return nil
}

func TestClientDoesNotReconnectAfterClose(t *testing.T) {
	session := &closeGuardSession{}
	connectCalls := 0
	client := NewClient(func(context.Context) (Session, error) {
		connectCalls++
		return session, nil
	})

	if _, err := client.Tools(context.Background()); err != nil {
		t.Fatalf("initial Tools: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	_, err := client.Tools(context.Background())
	if !errors.Is(err, ErrClientClosed) {
		t.Fatalf("Tools after Close error = %v, want ErrClientClosed", err)
	}
	if connectCalls != 1 {
		t.Fatalf("connect calls = %d, want 1", connectCalls)
	}
}
