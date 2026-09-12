package session

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadJournalRejectsSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows symlink creation may require elevated privileges")
	}

	root := t.TempDir()
	target := filepath.Join(root, "outside.json")
	if err := os.WriteFile(target, []byte("[]\n"), 0o600); err != nil {
		t.Fatalf("seed symlink target: %v", err)
	}

	path := filepath.Join(root, "session.json")
	if err := os.Symlink(target, path); err != nil {
		t.Fatalf("create journal symlink: %v", err)
	}

	if _, err := LoadJournal(path); err == nil {
		t.Fatal("LoadJournal() error = nil, want symlink rejection")
	}
}
