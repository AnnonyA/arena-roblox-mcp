package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvAcceptsUTF8BOM(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	data := append([]byte{0xEF, 0xBB, 0xBF}, []byte("ARENA_API_KEY=test-key\n")...)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	values := make(map[string]string)
	lookup := func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
	set := func(key, value string) error {
		values[key] = value
		return nil
	}

	if err := LoadDotEnv(path, lookup, set); err != nil {
		t.Fatalf("LoadDotEnv: %v", err)
	}
	if got := values["ARENA_API_KEY"]; got != "test-key" {
		t.Fatalf("ARENA_API_KEY = %q, want %q; loaded values: %#v", got, "test-key", values)
	}
}
