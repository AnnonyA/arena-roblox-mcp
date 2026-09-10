package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDotEnvRejectsNULBeforeApplyingVariables(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	data := []byte("ARENA_RBX_GOOD=value\nARENA_RBX_BAD=bad\x00value\n")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	setCalls := 0
	set := func(key, value string) error {
		setCalls++
		if strings.ContainsRune(key, '\x00') || strings.ContainsRune(value, '\x00') {
			return fmt.Errorf("environment variable contains NUL")
		}
		return nil
	}

	err := LoadDotEnv(path, func(string) (string, bool) { return "", false }, set)
	if err == nil {
		t.Fatal("LoadDotEnv error = nil, want NUL validation error")
	}
	if setCalls != 0 {
		t.Fatalf("set calls = %d, want 0 before validation completes", setCalls)
	}
}
