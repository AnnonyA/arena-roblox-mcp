package arena

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatRejectsOversizedToolCallID(t *testing.T) {
	t.Parallel()

	oversizedID := strings.Repeat("a", 4097)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = fmt.Fprintf(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":%q,\"type\":\"function\",\"function\":{\"name\":\"inspect\",\"arguments\":\"{}\"}}]}}]}\n\n", oversizedID)
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err == nil {
		t.Fatal("StreamChat error = nil, want oversized tool call id error")
	}
	if !strings.Contains(err.Error(), "tool call id exceeds") {
		t.Fatalf("StreamChat error = %q, want oversized tool call id error", err)
	}
}
