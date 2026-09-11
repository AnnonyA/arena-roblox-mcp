package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDotEnvAcceptsLineLargerThanScannerDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	value := strings.Repeat("x", 70<<10)
	if err := os.WriteFile(path, []byte("ARENA_API_KEY="+value+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	loaded := make(map[string]string)
	err := LoadDotEnv(path,
		func(key string) (string, bool) { return "", false },
		func(key, value string) error {
			loaded[key] = value
			return nil
		},
	)
	if err != nil {
		t.Fatalf("LoadDotEnv returned error for line below file limit: %v", err)
	}
	if got := loaded["ARENA_API_KEY"]; got != value {
		t.Fatalf("loaded value length = %d, want %d", len(got), len(value))
	}
}
