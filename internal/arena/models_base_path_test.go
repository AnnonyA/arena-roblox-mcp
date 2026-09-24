package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListModelsPreservesConfiguredBasePath(t *testing.T) {
	t.Parallel()

	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"model-a"}]}`))
	}))
	defer srv.Close()

	client := NewClient(ClientOptions{BaseURL: srv.URL + "/arena-proxy/"})
	models, err := client.ListModels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/arena-proxy/v1/models" {
		t.Fatalf("request path = %q, want %q", gotPath, "/arena-proxy/v1/models")
	}
	if len(models) != 1 || models[0].ID != "model-a" {
		t.Fatalf("models = %#v", models)
	}
}
