package arena

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatAcceptsSSEPayloadWithinEventLimit(t *testing.T) {
	const prefix = `{"choices":[{"index":0,"delta":{"content":"`
	const suffix = `"}}]}`
	payloadSize := maxSSEEventBytes - 1
	contentSize := payloadSize - len(prefix) - len(suffix)
	payload := prefix + strings.Repeat("x", contentSize) + suffix

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "data: %s\n\ndata: [DONE]\n\n", payload)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{BaseURL: server.URL, HTTPClient: server.Client()})
	result, err := client.StreamChat(context.Background(), ChatRequest{Model: "test"}, nil)
	if err != nil {
		t.Fatalf("StreamChat() error = %v", err)
	}
	if got, want := len(result.Text), contentSize; got != want {
		t.Fatalf("len(result.Text) = %d, want %d", got, want)
	}
}
