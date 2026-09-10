package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsBlankModelID(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"model-a"},{"id":"   "}]}`))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected blank model ID error")
	}
	if !strings.Contains(err.Error(), "blank model id") {
		t.Fatalf("error = %q, want blank model id error", err)
	}
}
