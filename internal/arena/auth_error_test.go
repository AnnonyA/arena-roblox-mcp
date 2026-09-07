package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArenaAuthenticationErrorsAreActionable(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := NewClient(ClientOptions{BaseURL: srv.URL, APIKey: "top-secret"})
	const want = "Arena authentication failed. Check your API key."

	t.Run("models", func(t *testing.T) {
		_, err := client.ListModels(context.Background())
		if err == nil {
			t.Fatal("expected authentication error")
		}
		if err.Error() != want {
			t.Fatalf("error = %q, want %q", err, want)
		}
	})

	t.Run("chat", func(t *testing.T) {
		_, err := client.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
		if err == nil {
			t.Fatal("expected authentication error")
		}
		if err.Error() != want {
			t.Fatalf("error = %q, want %q", err, want)
		}
	})
}
