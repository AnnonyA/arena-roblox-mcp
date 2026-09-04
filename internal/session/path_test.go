package session

import (
	"path/filepath"
	"testing"
)

func TestResolveDataDirUsesLocalAppData(t *testing.T) {
	got, err := resolveDataDir(`C:\\Users\\test\\AppData\\Local`, func() (string, error) {
		t.Fatal("fallback should not be called when LOCALAPPDATA is available")
		return "", nil
	})
	if err != nil {
		t.Fatalf("resolveDataDir() error = %v", err)
	}

	want := filepath.Join(`C:\\Users\\test\\AppData\\Local`, "arena-rbx", "sessions")
	if got != want {
		t.Fatalf("resolveDataDir() = %q, want %q", got, want)
	}
}

func TestResolveDataDirFallsBackToUserConfigDir(t *testing.T) {
	got, err := resolveDataDir("", func() (string, error) {
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
