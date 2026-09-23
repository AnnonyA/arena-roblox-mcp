package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatDoesNotFollowAuthenticatedRedirect(t *testing.T) {
	t.Parallel()

	const apiKey = "stream-redirect-secret"
	targetHit := false
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetHit = true
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("redirect target received Authorization header %q", got)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer target.Close()

	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), "Bearer "+apiKey; got != want {
			t.Errorf("source Authorization = %q, want %q", got, want)
		}
		http.Redirect(w, r, target.URL+"/v1/chat/completions", http.StatusTemporaryRedirect)
	}))
	defer source.Close()

	client := NewClient(ClientOptions{BaseURL: source.URL, APIKey: apiKey})
	_, err := client.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err == nil {
		t.Fatal("expected redirect to be rejected")
	}
	if targetHit {
		t.Fatal("authenticated Arena streaming request followed redirect")
	}
	if strings.Contains(err.Error(), apiKey) {
		t.Fatalf("error leaked API key: %q", err)
	}
}
