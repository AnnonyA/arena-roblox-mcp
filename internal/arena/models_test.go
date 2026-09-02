package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsUsesBearerAuthAndParsesIDs(t *testing.T) {
	t.Parallel()
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"model-a"},{"id":"model-b"}]}`))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL, APIKey: "secret-key"})
	got, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if auth != "Bearer secret-key" {
		t.Fatalf("Authorization = %q", auth)
	}
	if len(got) != 2 || got[0].ID != "model-a" || got[1].ID != "model-b" {
		t.Fatalf("models = %#v", got)
	}
}

func TestListModelsErrorDoesNotLeakAPIKey(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL, APIKey: "top-secret"})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "top-secret") {
		t.Fatalf("error leaked key: %v", err)
	}
	if !strings.Contains(err.Error(), "401") {
		t.Fatalf("error = %v", err)
	}
}
