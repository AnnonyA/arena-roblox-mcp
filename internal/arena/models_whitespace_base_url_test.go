package arena

import (
	"context"
	"strings"
	"testing"
)

func TestListModelsRejectsWhitespaceOnlyBaseURL(t *testing.T) {
	t.Parallel()

	c := NewClient(ClientOptions{BaseURL: " \t\r\n "})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); !strings.Contains(got, "base URL is empty") {
		t.Fatalf("error = %q, want empty base URL error", got)
	}
}
