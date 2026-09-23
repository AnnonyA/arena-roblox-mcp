package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatForbiddenErrorDoesNotLeakAPIKey(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL, APIKey: "top-secret"})
	_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "top-secret") {
		t.Fatalf("error leaked key: %v", err)
	}
	const want = "Arena authentication failed. Check your API key."
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err, want)
	}
}
