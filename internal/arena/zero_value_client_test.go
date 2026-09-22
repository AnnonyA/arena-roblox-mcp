package arena

import (
	"context"
	"strings"
	"testing"
)

func TestListModelsRejectsZeroValueClient(t *testing.T) {
	t.Parallel()

	client := &Client{}
	_, err := client.ListModels(context.Background())
	if err == nil || !strings.Contains(err.Error(), "HTTP client is nil") {
		t.Fatalf("zero-value Client.ListModels error = %v, want HTTP client is nil", err)
	}
}
