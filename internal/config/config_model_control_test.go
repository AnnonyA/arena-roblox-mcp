package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsConfiguredModelControlCharacter(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	if err := os.WriteFile(path, []byte(`{"arena":{"model":"model\nother"}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want configured model control-character error")
	}
	if !strings.Contains(err.Error(), "arena.model") || !strings.Contains(err.Error(), "control character") {
		t.Fatalf("Load() error = %q, want arena.model control-character error", err)
	}
}
