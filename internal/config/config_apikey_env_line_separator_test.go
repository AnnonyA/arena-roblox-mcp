package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsLineSeparatorsInAPIKeyEnv(t *testing.T) {
	for _, escaped := range []string{`\u2028`, `\u2029`} {
		t.Run(escaped, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "arena-rbx.json")
			data := []byte(`{"arena":{"apiKeyEnv":"CUSTOM` + escaped + `ARENA_KEY"}}`)
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatalf("WriteFile: %v", err)
			}

			_, err := Load(path)
			if err == nil {
				t.Fatal("Load succeeded, want line-separator validation error")
			}
			if !strings.Contains(err.Error(), "arena.apiKeyEnv must not contain line separators") {
				t.Fatalf("Load error = %q, want arena.apiKeyEnv line-separator error", err)
			}
		})
	}
}
