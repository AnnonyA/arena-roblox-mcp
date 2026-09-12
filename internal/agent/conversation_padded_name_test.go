package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestRunConversationRejectsPaddedToolCallNameBeforeDispatch(t *testing.T) {
	dispatcher := NewToolDispatcher([]string{"script_read"}, &conversationToolCaller{})
	initial := []arena.Message{{Role: "user", Content: "inspect safely"}}
	rounds := 0
	runRound := func(_ context.Context, _ []arena.Message) (arena.ChatResult, error) {
		rounds++
		return arena.ChatResult{ToolCalls: []arena.ToolCall{{
			ID:   "call-padded-name",
			Type: "function",
			Function: arena.FunctionCall{
				Name:      " script_read",
				Arguments: `{}`,
			},
		}}}, nil
	}

	_, err := RunConversation(context.Background(), 1, initial, runRound, dispatcher)
	if !errors.Is(err, ErrInvalidToolCall) {
		t.Fatalf("RunConversation error = %v, want ErrInvalidToolCall", err)
	}
	if rounds != 1 {
		t.Fatalf("chat rounds = %d, want 1", rounds)
	}
}
