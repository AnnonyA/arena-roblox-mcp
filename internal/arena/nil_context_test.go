package arena

import (
	"strings"
	"testing"
)

func TestListModelsRejectsNilContext(t *testing.T) {
	t.Parallel()

	client := NewClient(ClientOptions{BaseURL: "http://127.0.0.1"})
	_, err := client.ListModels(nil)
	if err == nil || !strings.Contains(err.Error(), "context is nil") {
		t.Fatalf("ListModels(nil) error = %v, want context is nil", err)
	}
}
