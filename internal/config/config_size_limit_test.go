package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsOversizedConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := `{"arena":{"model":"` + strings.Repeat("x", (1<<20)+1) + `"}}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load error = nil, want oversized config error")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Fatalf("Load error = %q, want size-limit error", err)
	}
}
