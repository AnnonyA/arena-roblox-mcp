package arena

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsTooManyModels(t *testing.T) {
	t.Parallel()

	var body strings.Builder
	body.WriteString(`{"data":[`)
	for i := 0; i < 4097; i++ {
		if i > 0 {
			body.WriteByte(',')
		}
		fmt.Fprintf(&body, `{"id":"model-%d"}`, i)
	}
	body.WriteString(`]}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body.String()))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected model count limit error")
	}
	if !strings.Contains(err.Error(), "model count exceeds") {
		t.Fatalf("error = %q, want model count limit error", err)
	}
}
