package arena

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func FuzzListModelsModelIDValidation(f *testing.F) {
	for _, seed := range []string{
		"safe-model",
		" model",
		"model\nname",
		"safe\u200bhidden",
		"safe\u202ehidden",
		"モデル-1",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, id string) {
		if !utf8.ValidString(id) || len(id) > maxModelIDBytes {
			t.Skip()
		}

		body, err := json.Marshal(modelsResponse{Data: []Model{{ID: id}}})
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(body)
		}))
		defer server.Close()

		models, err := NewClient(ClientOptions{BaseURL: server.URL}).ListModels(context.Background())
		if err != nil {
			return
		}
		if len(models) != 1 || models[0].ID != id {
			t.Fatalf("ListModels() = %#v, want exact model ID %q", models, id)
		}
		if strings.TrimSpace(id) != id || id == "" {
			t.Fatalf("ListModels() accepted blank or surrounding whitespace model ID %q", id)
		}
		for _, r := range id {
			if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || r == '\u2028' || r == '\u2029' {
				t.Fatalf("ListModels() accepted unsafe rune U+%04X in model ID %q", r, id)
			}
		}
	})
}
