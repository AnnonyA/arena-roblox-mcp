package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRejectsInvalidUTF8(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte{'{', '"', 'a', 'r', 'e', 'n', 'a', '"', ':', '{', '"', 'm', 'o', 'd', 'e', 'l', '"', ':', '"', 0xff, '"', '}', '}'}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("Load error = nil, want invalid UTF-8 error")
	}
}
