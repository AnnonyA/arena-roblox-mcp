package agent

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
	mcppkg "github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
)

type conversationToolCaller struct {
	calls []string
}

func (c *conversationToolCaller) CallTool(_ context.Context, name string, arguments json.RawMessage) (mcppkg.ToolResult, error) {
	c.calls = append(c.calls, name+":"+string(arguments))
	return mcppkg.ToolResult{StructuredContent: json.RawMessage(`{"source":"print('ok')"}`)}, nil
}

func TestRunConversationFeedsToolResultsBackToArena(t *testing.T) {
	caller := &conversationToolCaller{}
	dispatcher := NewToolDispatcher([]string{"script_read"}, caller)
	initial := []arena.Message{{Role: "user", Content: "inspect main script"}}

	var requests [][]arena.Message
	runRound := func(_ context.Context, messages []arena.Message) (arena.ChatResult, error) {
		requests = append(requests, append([]arena.Message(nil), messages...))
		if len(requests) == 1 {
			return arena.ChatResult{ToolCalls: []arena.ToolCall{{ID: "call-1", Type: "function", Function: arena.FunctionCall{Name: "script_read", Arguments: `{"path":"ServerScriptService.Main"}`}}}}, nil
		}
		return arena.ChatResult{Text: "inspection complete"}, nil
	}

	text, err := RunConversation(context.Background(), 12, initial, runRound, dispatcher)
	if err != nil {
		t.Fatalf("RunConversation: %v", err)
	}
	if text != "inspection complete" {
		t.Fatalf("text = %q, want inspection complete", text)
	}
	if got, want := caller.calls, []string{`script_read:{"path":"ServerScriptService.Main"}`}; !reflect.DeepEqual(got, want) {
		t.Fatalf("tool calls = %#v, want %#v", got, want)
	}
	if len(requests) != 2 {
		t.Fatalf("Arena requests = %d, want 2", len(requests))
	}
	second := requests[1]
	if len(second) != 3 {
		t.Fatalf("second request messages = %#v, want user + assistant tool call + tool result", second)
	}
	if second[1].Role != "assistant" || len(second[1].ToolCalls) != 1 || second[1].ToolCalls[0].ID != "call-1" {
		t.Fatalf("assistant tool-call message = %#v", second[1])
	}
	if second[2].Role != "tool" || second[2].ToolCallID != "call-1" || second[2].Content != `{"source":"print('ok')"}` {
		t.Fatalf("tool result message = %#v", second[2])
	}
}

func TestRunConversationFeedsToolErrorsBackToArena(t *testing.T) {
	dispatcher := NewToolDispatcher([]string{"script_read"}, &conversationToolCaller{})
	initial := []arena.Message{{Role: "user", Content: "inspect something"}}

	var requests [][]arena.Message
	runRound := func(_ context.Context, messages []arena.Message) (arena.ChatResult, error) {
		requests = append(requests, append([]arena.Message(nil), messages...))
		if len(requests) == 1 {
			return arena.ChatResult{ToolCalls: []arena.ToolCall{{ID: "call-bad", Type: "function", Function: arena.FunctionCall{Name: "missing_tool", Arguments: `{}`}}}}, nil
		}
		return arena.ChatResult{Text: "recovered after tool error"}, nil
	}

	text, err := RunConversation(context.Background(), 12, initial, runRound, dispatcher)
	if err != nil {
		t.Fatalf("RunConversation: %v", err)
	}
	if text != "recovered after tool error" {
		t.Fatalf("text = %q, want recovered after tool error", text)
	}
	if len(requests) != 2 {
		t.Fatalf("Arena requests = %d, want 2", len(requests))
	}
	second := requests[1]
	if len(second) != 3 {
		t.Fatalf("second request messages = %#v, want user + assistant tool call + tool error", second)
	}
	if got := second[2]; got.Role != "tool" || got.ToolCallID != "call-bad" || got.Content != `{"error":"unknown tool: missing_tool"}` {
		t.Fatalf("tool error message = %#v", got)
	}
}

func TestRunConversationAllowsFinalAnswerAfterLastToolRound(t *testing.T) {
	caller := &conversationToolCaller{}
	dispatcher := NewToolDispatcher([]string{"script_read"}, caller)
	initial := []arena.Message{{Role: "user", Content: "inspect once"}}

	var requests int
	runRound := func(_ context.Context, _ []arena.Message) (arena.ChatResult, error) {
		requests++
		if requests == 1 {
			return arena.ChatResult{ToolCalls: []arena.ToolCall{{ID: "call-1", Type: "function", Function: arena.FunctionCall{Name: "script_read", Arguments: `{}`}}}}, nil
		}
		return arena.ChatResult{Text: "done after one tool round"}, nil
	}

	text, err := RunConversation(context.Background(), 1, initial, runRound, dispatcher)
	if err != nil {
		t.Fatalf("RunConversation: %v", err)
	}
	if text != "done after one tool round" {
		t.Fatalf("text = %q, want final answer after the last permitted tool round", text)
	}
	if requests != 2 {
		t.Fatalf("Arena requests = %d, want 2", requests)
	}
}
