package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsInvalidUTF8(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte{'{', '"', 'd', 'a', 't', 'a', '"', ':', '[', '{', '"', 'i', 'd', '"', ':', '"', 0xff, '"', '}', ']', '}'})
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected invalid UTF-8 error")
	}
	if !strings.Contains(err.Error(), "invalid UTF-8") {
		t.Fatalf("error = %q, want invalid UTF-8 error", err)
	}
}
