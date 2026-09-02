package mcp

import (
	"context"
	"encoding/json"
	"testing"
)

type callSession struct {
	name      string
	arguments json.RawMessage
}

func (s *callSession) ListTools(context.Context) ([]Tool, error) {
	return []Tool{{Name: "read_script"}}, nil
}

func (s *callSession) CallTool(_ context.Context, name string, arguments json.RawMessage) (ToolResult, error) {
	s.name = name
	s.arguments = append(json.RawMessage(nil), arguments...)
	return ToolResult{Content: json.RawMessage(`[{"type":"text","text":"ok"}]`)}, nil
}

func (s *callSession) Close() error { return nil }

func TestClientCallToolReusesPersistentSession(t *testing.T) {
	session := &callSession{}
	connectCalls := 0
	client := NewClient(func(context.Context) (Session, error) {
		connectCalls++
		return session, nil
	})

	if _, err := client.Tools(context.Background()); err != nil {
		t.Fatalf("Tools: %v", err)
	}
	arguments := json.RawMessage(`{"script":"ServerScriptService.Main"}`)
	result, err := client.CallTool(context.Background(), "read_script", arguments)
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}

	if connectCalls != 1 {
		t.Fatalf("connect calls = %d, want 1", connectCalls)
	}
	if session.name != "read_script" {
		t.Fatalf("tool name = %q, want read_script", session.name)
	}
	if string(session.arguments) != string(arguments) {
		t.Fatalf("tool arguments = %s, want %s", session.arguments, arguments)
	}
	if string(result.Content) != `[{"type":"text","text":"ok"}]` {
		t.Fatalf("result content = %s", result.Content)
	}
}
