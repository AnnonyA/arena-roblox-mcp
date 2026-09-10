package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsBidirectionalFormattingInModelID(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"safe\u202emodel"}]}`))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected bidirectional-formatting model ID error")
	}
	if !strings.Contains(err.Error(), "bidirectional formatting") {
		t.Fatalf("error = %q, want bidirectional-formatting model ID error", err)
	}
}
