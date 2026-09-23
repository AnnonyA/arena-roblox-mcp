package arena

import (
	"context"
	"strings"
	"testing"
)

func TestListModelsRejectsBaseURLWithoutHTTPScheme(t *testing.T) {
	client := NewClient(ClientOptions{BaseURL: "api.arena.ai"})

	_, err := client.ListModels(context.Background())
	if err == nil {
		t.Fatal("ListModels() error = nil, want invalid base URL error")
	}
	if !strings.Contains(err.Error(), "base URL must use http or https") {
		t.Fatalf("ListModels() error = %q, want http/https scheme error", err)
	}
}
