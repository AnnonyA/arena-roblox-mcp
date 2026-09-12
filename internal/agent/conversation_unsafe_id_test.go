package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestRunConversationRejectsUnsafeToolCallIDBeforeDispatch(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "control character", id: "call-\x1b1"},
		{name: "unicode format character", id: "call-\u200b1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caller := &conversationToolCaller{}
			dispatcher := NewToolDispatcher([]string{"script_read"}, caller)
			initial := []arena.Message{{Role: "user", Content: "inspect safely"}}
			runRound := func(_ context.Context, _ []arena.Message) (arena.ChatResult, error) {
				return arena.ChatResult{ToolCalls: []arena.ToolCall{{
					ID:   tt.id,
					Type: "function",
					Function: arena.FunctionCall{
						Name:      "script_read",
						Arguments: `{}`,
					},
				}}}, nil
			}

			_, err := RunConversation(context.Background(), 1, initial, runRound, dispatcher)
			if !errors.Is(err, ErrInvalidToolCall) {
				t.Fatalf("RunConversation error = %v, want ErrInvalidToolCall", err)
			}
			if len(caller.calls) != 0 {
				t.Fatalf("backend calls = %#v, want none", caller.calls)
			}
		})
	}
}
