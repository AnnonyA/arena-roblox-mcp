package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsOversizedModelID(t *testing.T) {
	t.Parallel()

	oversizedID := strings.Repeat("a", 4097)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"` + oversizedID + `"}]}`))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected oversized model ID error")
	}
	if !strings.Contains(err.Error(), "model id exceeds") {
		t.Fatalf("error = %q, want oversized model ID error", err)
	}
}
