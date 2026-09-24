package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestListModelsDoesNotFollowRedirects(t *testing.T) {
	t.Parallel()

	var targetRequests atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetRequests.Add(1)
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("redirect target received Authorization header %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"unexpected"}]}`))
	}))
	defer target.Close()

	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret-key" {
			t.Fatalf("source Authorization = %q", got)
		}
		http.Redirect(w, r, target.URL+"/v1/models", http.StatusFound)
	}))
	defer source.Close()

	client := NewClient(ClientOptions{BaseURL: source.URL, APIKey: "secret-key"})
	_, err := client.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected redirect response to fail")
	}
	const want = "Arena models request failed with HTTP 302"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err, want)
	}
	if got := targetRequests.Load(); got != 0 {
		t.Fatalf("redirect target requests = %d, want 0", got)
	}
}
