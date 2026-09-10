package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAcceptsUTF8BOM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := append([]byte{0xEF, 0xBB, 0xBF}, []byte(`{"arena":{"model":"bom-model"}}`)...)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load error = %v, want nil", err)
	}
	if cfg.Arena.Model != "bom-model" {
		t.Fatalf("Arena.Model = %q, want %q", cfg.Arena.Model, "bom-model")
	}
}
