package session

import (
	"path/filepath"
	"testing"
)

func TestResolveDataDirFallsBackWhenLocalAppDataIsBlank(t *testing.T) {
	got, err := resolveDataDir("   \t", func() (string, error) {
		return filepath.Join("home", "test", ".config"), nil
	})
	if err != nil {
		t.Fatalf("resolveDataDir() error = %v", err)
	}

	want := filepath.Join("home", "test", ".config", "arena-rbx", "sessions")
	if got != want {
		t.Fatalf("resolveDataDir() = %q, want %q", got, want)
	}
}
