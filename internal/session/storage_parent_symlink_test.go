package session

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSaveJournalRejectsSymlinkedSessionDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows symlink creation may require elevated privileges")
	}

	root := t.TempDir()
	outside := filepath.Join(root, "outside")
	if err := os.Mkdir(outside, 0o700); err != nil {
		t.Fatalf("create outside directory: %v", err)
	}

	sessions := filepath.Join(root, "sessions")
	if err := os.Symlink(outside, sessions); err != nil {
		t.Fatalf("create session directory symlink: %v", err)
	}

	path := filepath.Join(sessions, "session.json")
	err := SaveJournal(path, NewJournal())
	if err == nil {
		t.Fatal("SaveJournal() error = nil, want symlinked directory rejection")
	}
	if !strings.Contains(err.Error(), "symbolic links") {
		t.Fatalf("SaveJournal() error = %q, want symbolic-link rejection", err)
	}
	if _, statErr := os.Stat(filepath.Join(outside, "session.json")); !os.IsNotExist(statErr) {
		t.Fatalf("outside journal unexpectedly created: %v", statErr)
	}
}
