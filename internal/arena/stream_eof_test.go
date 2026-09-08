package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStreamChatDispatchesFinalSSEEventAtEOF(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"ok\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	got, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "ok" {
		t.Fatalf("text = %q, want ok", got.Text)
	}
}
