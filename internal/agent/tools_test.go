package agent

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	mcppkg "github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
)

type recordingToolCaller struct {
	calls int
	name  string
	args  json.RawMessage
}

func (c *recordingToolCaller) CallTool(_ context.Context, name string, args json.RawMessage) (mcppkg.ToolResult, error) {
	c.calls++
	c.name = name
	c.args = append(json.RawMessage(nil), args...)
	return mcppkg.ToolResult{StructuredContent: json.RawMessage(`{"ok":true}`)}, nil
}

func TestToolDispatcherRejectsUnknownToolWithoutCallingBackend(t *testing.T) {
	caller := &recordingToolCaller{}
	dispatcher := NewToolDispatcher([]string{"read_script"}, caller)

	_, err := dispatcher.Dispatch(context.Background(), "delete_everything", json.RawMessage(`{}`))

	if !errors.Is(err, ErrUnknownTool) {
		t.Fatalf("Dispatch() error = %v, want ErrUnknownTool", err)
	}
	if caller.calls != 0 {
		t.Fatalf("backend calls = %d, want 0", caller.calls)
	}
}

func TestToolDispatcherCallsKnownTool(t *testing.T) {
	caller := &recordingToolCaller{}
	dispatcher := NewToolDispatcher([]string{"read_script"}, caller)
	args := json.RawMessage(`{"path":"ServerScriptService.Main"}`)

	result, err := dispatcher.Dispatch(context.Background(), "read_script", args)
	if err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}
	if caller.calls != 1 {
		t.Fatalf("backend calls = %d, want 1", caller.calls)
	}
	if caller.name != "read_script" {
		t.Fatalf("backend tool = %q, want read_script", caller.name)
	}
	if string(caller.args) != string(args) {
		t.Fatalf("backend args = %s, want %s", caller.args, args)
	}
	if string(result.StructuredContent) != `{"ok":true}` {
		t.Fatalf("result = %s, want backend result", result.StructuredContent)
	}
}
