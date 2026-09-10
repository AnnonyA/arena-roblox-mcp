package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRejectsNullConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	if err := os.WriteFile(path, []byte("null\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("Load error = nil, want null config error")
	}
}
