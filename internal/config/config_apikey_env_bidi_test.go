package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsBidiFormattingInAPIKeyEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"arena":{"apiKeyEnv":"CUSTOM\u202eARENA_KEY"}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load succeeded, want bidirectional-formatting validation error")
	}
	if !strings.Contains(err.Error(), "arena.apiKeyEnv must not contain bidirectional formatting") {
		t.Fatalf("Load error = %q, want arena.apiKeyEnv bidirectional-formatting error", err)
	}
}
