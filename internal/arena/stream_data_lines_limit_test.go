package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatRejectsTooManySSEDataLines(t *testing.T) {
	var body strings.Builder
	for i := 0; i < 4097; i++ {
		body.WriteString("data:\n")
	}
	body.WriteString("\n")
	body.WriteString("data: [DONE]\n\n")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(body.String()))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", server.Client())
	_, err := client.StreamChat(context.Background(), ChatRequest{Model: "test"}, nil)
	if err == nil || !strings.Contains(err.Error(), "too many data lines") {
		t.Fatalf("StreamChat error = %v, want too many data lines error", err)
	}
}
