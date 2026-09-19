package arena

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListModelsRejectsUnicodeFormatCharacterInID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rune rune
	}{
		{name: "soft hyphen", rune: '\u00ad'},
		{name: "zero width space", rune: '\u200b'},
		{name: "word joiner", rune: '\u2060'},
		{name: "byte order mark", rune: '\ufeff'},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"data":[{"id":"safe%chidden"}]}`, tt.rune)
			}))
			defer server.Close()

			client := NewClient(ClientOptions{BaseURL: server.URL})
			_, err := client.ListModels(context.Background())
			if err == nil || !strings.Contains(err.Error(), "format character") {
				t.Fatalf("ListModels() error = %v, want Unicode format character rejection for U+%04X", err, tt.rune)
			}
		})
	}
}
