package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatRejectsDuplicateToolCallIDsAcrossIndexes(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call_dup\",\"type\":\"function\",\"function\":{\"name\":\"script_read\",\"arguments\":\"{}\"}},{\"index\":1,\"id\":\"call_dup\",\"type\":\"function\",\"function\":{\"name\":\"script_search\",\"arguments\":\"{}\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err == nil {
		t.Fatal("StreamChat error = nil, want duplicate tool-call id error")
	}
	if !strings.Contains(err.Error(), "duplicate tool call id") {
		t.Fatalf("StreamChat error = %q, want duplicate tool-call id error", err)
	}
}
