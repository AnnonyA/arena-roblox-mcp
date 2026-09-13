package agent

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestRunConversationRejectsOversizedToolCallMetadataBeforeDispatch(t *testing.T) {
	tests := []struct {
		name string
		call arena.ToolCall
	}{
		{
			name: "id",
			call: arena.ToolCall{
				ID:   strings.Repeat("a", 257),
				Type: "function",
				Function: arena.FunctionCall{
					Name:      "script_read",
					Arguments: `{}`,
				},
			},
		},
		{
			name: "name",
			call: arena.ToolCall{
				ID:   "call-oversized-name",
				Type: "function",
				Function: arena.FunctionCall{
					Name:      strings.Repeat("a", 257),
					Arguments: `{}`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caller := &conversationToolCaller{}
			dispatcher := NewToolDispatcher([]string{"script_read"}, caller)
			initial := []arena.Message{{Role: "user", Content: "inspect safely"}}
			runRound := func(_ context.Context, _ []arena.Message) (arena.ChatResult, error) {
				return arena.ChatResult{ToolCalls: []arena.ToolCall{tt.call}}, nil
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
