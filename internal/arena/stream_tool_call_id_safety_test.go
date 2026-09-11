package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatRejectsUnsafeToolCallIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   string
		want string
	}{
		{name: "control character", id: "call\\n1", want: "control character"},
		{name: "bidirectional formatting", id: "call\\u202e1", want: "bidirectional formatting"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"" + tt.id + "\",\"type\":\"function\",\"function\":{\"name\":\"script_read\",\"arguments\":\"{}\"}}]}}]}\n\n"))
				_, _ = w.Write([]byte("data: [DONE]\n\n"))
			}))
			defer srv.Close()

			c := NewClient(ClientOptions{BaseURL: srv.URL})
			_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
			if err == nil {
				t.Fatalf("StreamChat error = nil, want %s error", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("StreamChat error = %q, want %s error", err, tt.want)
			}
		})
	}
}
