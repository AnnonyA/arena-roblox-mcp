package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type cancelDuringCallSession struct {
	started chan struct{}
	release chan struct{}
}

func (s *cancelDuringCallSession) ListTools(context.Context) ([]Tool, error) { return nil, nil }

func (s *cancelDuringCallSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	close(s.started)
	<-s.release
	return ToolResult{Content: json.RawMessage(`{"ok":true}`)}, nil
}

func (s *cancelDuringCallSession) Close() error { return nil }

func TestClientCallToolDiscardsResultCanceledDuringCall(t *testing.T) {
	session := &cancelDuringCallSession{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	client := NewClient(func(context.Context) (Session, error) { return session, nil })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		_, err := client.CallTool(ctx, "read_script", nil)
		errCh <- err
	}()

	<-session.started
	cancel()
	close(session.release)

	if err := <-errCh; !errors.Is(err, context.Canceled) {
		t.Fatalf("CallTool error = %v, want context.Canceled", err)
	}
}
