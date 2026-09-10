package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsOversizedResponse(t *testing.T) {
	t.Parallel()

	const responseLimit = 1024 * 1024
	largeID := strings.Repeat("x", responseLimit)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"` + largeID + `"}]}`))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("ListModels error = nil, want oversized response error")
	}
	if !strings.Contains(err.Error(), "models response exceeds") {
		t.Fatalf("error = %q, want models response limit error", err)
	}
}
