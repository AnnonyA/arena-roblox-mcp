package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
	mcppkg "github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
)

func TestRunConversationRejectsNonObjectToolArgumentsBeforeDispatch(t *testing.T) {
	called := false
	dispatcher := NewToolDispatcher([]string{"script_read"}, toolCallerFunc(func(context.Context, string, []byte) (mcppkg.ToolResult, error) {
		called = true
		return mcppkg.ToolResult{}, nil
	}))

	_, err := RunConversation(context.Background(), 1, nil, func(context.Context, []arena.Message) (arena.ChatResult, error) {
		return arena.ChatResult{ToolCalls: []arena.ToolCall{{
			ID: "call-1",
			Function: arena.ToolCallFunction{
				Name:      "script_read",
				Arguments: "null",
			},
		}}}, nil
	}, dispatcher)
	if !errors.Is(err, ErrInvalidToolCall) {
		t.Fatalf("RunConversation error = %v, want ErrInvalidToolCall", err)
	}
	if called {
		t.Fatal("malformed tool arguments reached the backend")
	}
}
