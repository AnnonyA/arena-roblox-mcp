package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsControlCharactersInAPIKeyEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"arena":{"apiKeyEnv":"CUSTOM\nARENA_KEY"}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load succeeded, want control-character validation error")
	}
	if !strings.Contains(err.Error(), "arena.apiKeyEnv must not contain control characters") {
		t.Fatalf("Load error = %q, want arena.apiKeyEnv control-character error", err)
	}
}
