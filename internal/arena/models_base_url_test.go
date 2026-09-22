package arena

import (
	"context"
	"strings"
	"testing"
)

func TestListModelsRejectsEmptyBaseURL(t *testing.T) {
	t.Parallel()

	c := NewClient(ClientOptions{})
	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); !strings.Contains(got, "base URL is empty") {
		t.Fatalf("error = %q, want empty base URL error", got)
	}
}
