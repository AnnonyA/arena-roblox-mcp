package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDotEnvRejectsInvalidUTF8(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte{'A', 'R', 'E', 'N', 'A', '_', 'A', 'P', 'I', '_', 'K', 'E', 'Y', '=', 0xff, '\n'}, 0o600); err != nil {
		t.Fatal(err)
	}

	setCalled := false
	err := LoadDotEnv(path, func(string) (string, bool) { return "", false }, func(string, string) error {
		setCalled = true
		return nil
	})
	if err == nil {
		t.Fatal("LoadDotEnv error = nil, want invalid UTF-8 error")
	}
	if !strings.Contains(err.Error(), "invalid UTF-8") {
		t.Fatalf("LoadDotEnv error = %q, want invalid UTF-8 error", err)
	}
	if setCalled {
		t.Fatal("setter called for invalid UTF-8 .env")
	}
}
