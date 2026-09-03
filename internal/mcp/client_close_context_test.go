package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type closeContextSession struct{ closed chan struct{} }

func (s *closeContextSession) ListTools(context.Context) ([]Tool, error) { return nil, nil }
func (s *closeContextSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{}, nil
}
func (s *closeContextSession) Close() error { close(s.closed); return nil }

func TestClientCloseContextCanCancelWhileConnectionIsInFlight(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	session := &closeContextSession{closed: make(chan struct{})}
	client := NewClient(func(context.Context) (Session, error) {
		close(started)
		<-release
		return session, nil
	})

	callDone := make(chan error, 1)
	go func() {
		_, err := client.CallTool(context.Background(), "read_script", nil)
		callDone <- err
	}()
	<-started

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err := client.CloseContext(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("CloseContext error = %v, want context.Canceled", err)
	}
	if time.Since(start) > 100*time.Millisecond {
		t.Fatal("CloseContext blocked on in-flight connection after cancellation")
	}

	close(release)
	if err := <-callDone; err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	select {
	case <-session.closed:
	default:
		t.Fatal("Close did not close connected session")
	}
}
