package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTrimsArenaFallbackModelWhitespace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"arena":{"fallbacks":["  fallback-model-id  "]}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := cfg.Arena.Fallbacks[0]; got != "fallback-model-id" {
		t.Fatalf("Arena.Fallbacks[0] = %q, want fallback-model-id", got)
	}
}
