package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatRejectsToolCallIDsWithSurroundingWhitespace(t *testing.T) {
	t.Parallel()

	for _, id := range []string{" call-1", "call-1 "} {
		id := id
		t.Run(id, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"" + id + "\",\"type\":\"function\",\"function\":{\"name\":\"script_read\",\"arguments\":\"{}\"}}]}}]}\n\n"))
				_, _ = w.Write([]byte("data: [DONE]\n\n"))
			}))
			defer srv.Close()

			c := NewClient(ClientOptions{BaseURL: srv.URL})
			_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
			if err == nil {
				t.Fatal("StreamChat error = nil, want surrounding whitespace error")
			}
			if !strings.Contains(err.Error(), "surrounding whitespace") {
				t.Fatalf("StreamChat error = %q, want surrounding whitespace error", err)
			}
		})
	}
}
