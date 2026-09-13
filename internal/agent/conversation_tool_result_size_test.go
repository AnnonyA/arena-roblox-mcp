package agent

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
	mcppkg "github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
)

type largeResultToolCaller struct {
	content json.RawMessage
}

func (c *largeResultToolCaller) CallTool(_ context.Context, _ string, _ json.RawMessage) (mcppkg.ToolResult, error) {
	return mcppkg.ToolResult{StructuredContent: c.content}, nil
}

func TestRunConversationCompactsOversizedToolResultBeforeNextArenaRound(t *testing.T) {
	payload, err := json.Marshal(strings.Repeat("é", 300<<10))
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	caller := &largeResultToolCaller{content: payload}
	dispatcher := NewToolDispatcher([]string{"script_read"}, caller)
	initial := []arena.Message{{Role: "user", Content: "inspect a large script"}}

	var requests [][]arena.Message
	runRound := func(_ context.Context, messages []arena.Message) (arena.ChatResult, error) {
		requests = append(requests, append([]arena.Message(nil), messages...))
		if len(requests) == 1 {
			return arena.ChatResult{ToolCalls: []arena.ToolCall{{ID: "call-large", Type: "function", Function: arena.FunctionCall{Name: "script_read", Arguments: `{}`}}}}, nil
		}
		return arena.ChatResult{Text: "done"}, nil
	}

	if _, err := RunConversation(context.Background(), 12, initial, runRound, dispatcher); err != nil {
		t.Fatalf("RunConversation: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("Arena requests = %d, want 2", len(requests))
	}
	content := requests[1][2].Content
	if len(content) > (256<<10)+64 {
		t.Fatalf("tool result bytes = %d, want bounded near 256 KiB", len(content))
	}
	if !strings.HasSuffix(content, "… [tool result truncated]") {
		t.Fatalf("tool result missing truncation marker")
	}
	if !utf8.ValidString(content) {
		t.Fatalf("compacted tool result is not valid UTF-8")
	}
}

func TestCompactToolResultSanitizesInvalidUTF8(t *testing.T) {
	content := compactToolResult(json.RawMessage{'o', 'k', 0xff})

	if !strings.HasPrefix(content, "ok") {
		t.Fatalf("compacted tool result = %q, want preserved valid prefix", content)
	}
	if !utf8.ValidString(content) {
		t.Fatalf("compacted tool result is not valid UTF-8")
	}
	if len(content) > maxToolResultBytes {
		t.Fatalf("tool result bytes = %d, want at most %d", len(content), maxToolResultBytes)
	}
}
