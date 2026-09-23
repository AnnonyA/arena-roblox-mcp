package arena

import (
	"context"
	"strings"
	"testing"
)

func TestListModelsRejectsBaseURLWithQueryOrFragment(t *testing.T) {
	t.Parallel()

	for _, baseURL := range []string{
		"https://api.arena.ai?tenant=test",
		"https://api.arena.ai#models",
	} {
		baseURL := baseURL
		t.Run(baseURL, func(t *testing.T) {
			t.Parallel()
			c := NewClient(ClientOptions{BaseURL: baseURL})
			_, err := c.ListModels(context.Background())
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), "base URL must not contain query or fragment") {
				t.Fatalf("error = %q", err)
			}
		})
	}
}
