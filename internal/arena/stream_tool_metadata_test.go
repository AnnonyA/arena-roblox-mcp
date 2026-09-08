package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStreamChatDoesNotDuplicateRepeatedToolMetadata(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call_1\",\"type\":\"function\",\"function\":{\"name\":\"script_read\",\"arguments\":\"{\\\"path\\\":\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call_1\",\"type\":\"function\",\"function\":{\"arguments\":\"\\\"A\\\"}\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	got, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.ToolCalls) != 1 {
		t.Fatalf("tool calls = %#v", got.ToolCalls)
	}
	call := got.ToolCalls[0]
	if call.ID != "call_1" {
		t.Fatalf("tool call id = %q, want call_1", call.ID)
	}
	if call.Type != "function" {
		t.Fatalf("tool call type = %q, want function", call.Type)
	}
	if call.Function.Name != "script_read" {
		t.Fatalf("tool function name = %q, want script_read", call.Function.Name)
	}
	if call.Function.Arguments != `{"path":"A"}` {
		t.Fatalf("tool arguments = %q, want %q", call.Function.Arguments, `{"path":"A"}`)
	}
}
