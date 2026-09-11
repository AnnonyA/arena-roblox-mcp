package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsDuplicateJSONKeys(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"safe"}],"data":[{"id":"override"}]}`))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("ListModels error = nil, want duplicate JSON key error")
	}
	if !strings.Contains(err.Error(), "duplicate JSON key") {
		t.Fatalf("ListModels error = %q, want duplicate JSON key error", err)
	}
}
