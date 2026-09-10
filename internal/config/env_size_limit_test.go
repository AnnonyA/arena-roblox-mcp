package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDotEnvRejectsOversizedFileBeforeSettingVariables(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	data := bytes.Repeat([]byte("A=value\n"), (1<<20)/len("A=value\n")+2)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	setCalled := false
	err := LoadDotEnv(path, func(string) (string, bool) { return "", false }, func(string, string) error {
		setCalled = true
		return nil
	})
	if err == nil {
		t.Fatal("LoadDotEnv error = nil, want oversized .env error")
	}
	if !strings.Contains(err.Error(), ".env file too large") {
		t.Fatalf("LoadDotEnv error = %q, want size-limit error", err)
	}
	if setCalled {
		t.Fatal("setter called before oversized .env was rejected")
	}
}
