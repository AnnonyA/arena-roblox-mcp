package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRejectsEmptyArenaFallbackModel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"arena":{"fallbacks":["   "]}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("Load succeeded with an empty Arena fallback model, want error")
	}
}
