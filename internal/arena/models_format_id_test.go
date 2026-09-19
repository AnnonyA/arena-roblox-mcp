package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsUnicodeFormatCharacterInID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"safe\u200bhidden"}]}`))
	}))
	defer server.Close()

	client := NewClient(ClientOptions{BaseURL: server.URL})
	_, err := client.ListModels(context.Background())
	if err == nil || !strings.Contains(err.Error(), "format character") {
		t.Fatalf("ListModels() error = %v, want Unicode format character rejection", err)
	}
}
