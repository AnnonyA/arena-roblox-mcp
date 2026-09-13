package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTrimsArenaModelWhitespace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"arena":{"model":"  arena-model-id  "}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Arena.Model != "arena-model-id" {
		t.Fatalf("Arena.Model = %q, want arena-model-id", cfg.Arena.Model)
	}
}
