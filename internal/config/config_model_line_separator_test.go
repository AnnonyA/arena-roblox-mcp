package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsLineSeparatorInArenaModel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte("{\"arena\":{\"model\":\"safe\\u2028model\"}}")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load succeeded, want line-separator validation error")
	}
	if !strings.Contains(err.Error(), "arena.model must not contain line separators") {
		t.Fatalf("Load error = %q, want arena.model line-separator error", err)
	}
}
