package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStreamChatPreservesArgumentDeltaThatPrefixesAccumulatedArguments(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call_1\",\"type\":\"function\",\"function\":{\"name\":\"script_read\",\"arguments\":\"{\\\"outer\\\":\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"function\":{\"arguments\":\"{\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"function\":{\"arguments\":\"}\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"function\":{\"arguments\":\"}\"}}]}}]}\n\n"))
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
	if got.ToolCalls[0].Function.Arguments != `{"outer":{}}` {
		t.Fatalf("tool arguments = %q, want %q", got.ToolCalls[0].Function.Arguments, `{"outer":{}}`)
	}
}
