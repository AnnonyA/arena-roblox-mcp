package arena

import (
	"context"
	"strings"
	"testing"
)

func TestListModelsRejectsNilClient(t *testing.T) {
	t.Parallel()

	var client *Client
	_, err := client.ListModels(context.Background())
	if err == nil || !strings.Contains(err.Error(), "client is nil") {
		t.Fatalf("nil Client.ListModels error = %v, want client is nil", err)
	}
}
