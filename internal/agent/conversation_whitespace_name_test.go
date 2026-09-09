package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestRunConversationRejectsWhitespaceToolCallNameBeforeDispatch(t *testing.T) {
	caller := &conversationToolCaller{}
	dispatcher := NewToolDispatcher([]string{"script_read"}, caller)
	initial := []arena.Message{{Role: "user", Content: "inspect safely"}}
	runRound := func(_ context.Context, _ []arena.Message) (arena.ChatResult, error) {
		return arena.ChatResult{ToolCalls: []arena.ToolCall{{
			ID:   "call-whitespace-name",
			Type: "function",
			Function: arena.FunctionCall{
				Name:      "   \t",
				Arguments: `{}`,
			},
		}}}, nil
	}

	_, err := RunConversation(context.Background(), 12, initial, runRound, dispatcher)
	if !errors.Is(err, ErrInvalidToolCall) {
		t.Fatalf("RunConversation error = %v, want ErrInvalidToolCall", err)
	}
	if len(caller.calls) != 0 {
		t.Fatalf("backend calls = %#v, want none", caller.calls)
	}
}
