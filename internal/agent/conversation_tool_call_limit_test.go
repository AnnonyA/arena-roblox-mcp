package agent

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestRunConversationRejectsTooManyToolCallsInSingleRoundBeforeDispatch(t *testing.T) {
	caller := &conversationToolCaller{}
	dispatcher := NewToolDispatcher([]string{"script_read"}, caller)
	initial := []arena.Message{{Role: "user", Content: "inspect safely"}}

	calls := make([]arena.ToolCall, maxToolCallsPerRound+1)
	for i := range calls {
		calls[i] = arena.ToolCall{
			ID:   fmt.Sprintf("call-%d", i),
			Type: "function",
			Function: arena.FunctionCall{
				Name:      "script_read",
				Arguments: `{}`,
			},
		}
	}

	runRound := func(_ context.Context, _ []arena.Message) (arena.ChatResult, error) {
		return arena.ChatResult{ToolCalls: calls}, nil
	}

	_, err := RunConversation(context.Background(), 12, initial, runRound, dispatcher)
	if !errors.Is(err, ErrInvalidToolCall) {
		t.Fatalf("RunConversation error = %v, want ErrInvalidToolCall", err)
	}
	if len(caller.calls) != 0 {
		t.Fatalf("backend calls = %d, want none", len(caller.calls))
	}
}
