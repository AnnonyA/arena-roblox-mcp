package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsDoesNotFollowAuthenticatedRedirect(t *testing.T) {
	t.Parallel()

	const apiKey = "redirect-secret"
	targetHit := false
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetHit = true
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("redirect target received Authorization header %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"model-a"}]}`))
	}))
	defer target.Close()

	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), "Bearer "+apiKey; got != want {
			t.Errorf("source Authorization = %q, want %q", got, want)
		}
		http.Redirect(w, r, target.URL+"/v1/models", http.StatusTemporaryRedirect)
	}))
	defer source.Close()

	client := NewClient(ClientOptions{BaseURL: source.URL, APIKey: apiKey})
	_, err := client.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected redirect to be rejected")
	}
	if targetHit {
		t.Fatal("authenticated Arena request followed redirect")
	}
	if strings.Contains(err.Error(), apiKey) {
		t.Fatalf("error leaked API key: %q", err)
	}
}
