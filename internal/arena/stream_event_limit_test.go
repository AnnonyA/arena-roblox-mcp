package arena

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStreamChatRejectsOversizedMultilineSSEEvent(t *testing.T) {
	t.Parallel()

	var body strings.Builder
	chunk := strings.Repeat("x", 64*1024)
	for i := 0; i < 17; i++ {
		body.WriteString("data: ")
		body.WriteString(chunk)
		body.WriteString("\n")
	}
	body.WriteString("\n")

	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body.String())),
			Request:    req,
		}, nil
	})}

	c := NewClient(ClientOptions{BaseURL: "https://arena.invalid", HTTPClient: client})
	_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err == nil {
		t.Fatal("expected oversized SSE event error")
	}
	if !strings.Contains(err.Error(), "Arena SSE event exceeds 1048576 bytes") {
		t.Fatalf("error = %q, want oversized SSE event error", err)
	}
}
