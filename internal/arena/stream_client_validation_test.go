package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsForbiddenDoesNotLeakAPIKey(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	const apiKey = "models-super-secret"
	client := NewClient(ClientOptions{BaseURL: srv.URL, APIKey: apiKey})
	_, err := client.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected authentication error")
	}
	if got, want := err.Error(), "Arena authentication failed. Check your API key."; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if strings.Contains(err.Error(), apiKey) {
		t.Fatalf("authentication error leaked API key: %q", err)
	}
}
