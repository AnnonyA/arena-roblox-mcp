package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsControlCharactersInArenaFallbackModel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	if err := os.WriteFile(path, []byte(`{"arena":{"fallbacks":["fallback\nmodel"]}}`), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load error = nil, want arena.fallbacks control-character error")
	}
	if !strings.Contains(err.Error(), "arena.fallbacks[0]") || !strings.Contains(err.Error(), "control characters") {
		t.Fatalf("Load error = %q, want arena.fallbacks[0] control-character error", err)
	}
}
