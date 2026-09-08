package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestRunConversationRejectsNonObjectToolArgumentsBeforeDispatch(t *testing.T) {
	caller := &conversationToolCaller{}
	dispatcher := NewToolDispatcher([]string{"script_read"}, caller)

	_, err := RunConversation(context.Background(), 1, nil, func(context.Context, []arena.Message) (arena.ChatResult, error) {
		return arena.ChatResult{ToolCalls: []arena.ToolCall{{
			ID:       "call-1",
			Function: arena.FunctionCall{Name: "script_read", Arguments: "null"},
		}}}, nil
	}, dispatcher)
	if !errors.Is(err, ErrInvalidToolCall) {
		t.Fatalf("RunConversation error = %v, want ErrInvalidToolCall", err)
	}
	if len(caller.calls) != 0 {
		t.Fatalf("backend calls = %#v, want none", caller.calls)
	}
}
