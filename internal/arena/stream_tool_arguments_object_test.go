package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatRejectsNonObjectToolCallArguments(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call-1\",\"type\":\"function\",\"function\":{\"name\":\"inspect\",\"arguments\":\"[]\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err == nil {
		t.Fatal("StreamChat error = nil, want non-object tool arguments error")
	}
	if !strings.Contains(err.Error(), "JSON object") {
		t.Fatalf("StreamChat error = %q, want JSON object error", err)
	}
}
