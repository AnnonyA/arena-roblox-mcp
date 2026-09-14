package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsUnicodeFormatCharactersInArenaModel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"arena":{"model":"safe\u200bmodel"}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load succeeded, want Unicode-format validation error")
	}
	if !strings.Contains(err.Error(), "arena.model must not contain Unicode format characters") {
		t.Fatalf("Load error = %q, want arena.model Unicode-format error", err)
	}
}
