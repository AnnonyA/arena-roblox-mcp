package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatKeepsPrimaryChoiceIsolated(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"primary\"}},{\"index\":1,\"delta\":{\"content\":\"secondary\",\"tool_calls\":[{\"index\":0,\"id\":\"call_other\",\"type\":\"function\",\"function\":{\"name\":\"other_tool\",\"arguments\":\"{}\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	got, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "primary" {
		t.Fatalf("text = %q, want primary", got.Text)
	}
	if len(got.ToolCalls) != 0 {
		t.Fatalf("tool calls = %#v, want none from secondary choice", got.ToolCalls)
	}
}

func TestStreamChatRejectsNegativeChoiceIndex(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":-1,\"delta\":{\"content\":\"invalid\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err == nil || !strings.Contains(err.Error(), "negative choice index") {
		t.Fatalf("StreamChat error = %v, want negative choice index error", err)
	}
}
