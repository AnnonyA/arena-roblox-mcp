package session

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLoadJournalRejectsSymlinkedSessionDirectoryAncestor(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows symlink creation may require elevated privileges")
	}

	root := t.TempDir()
	outside := filepath.Join(root, "outside")
	if err := os.Mkdir(outside, 0o700); err != nil {
		t.Fatalf("create outside directory: %v", err)
	}

	journalPath := filepath.Join(outside, "session.json")
	if err := SaveJournal(journalPath, NewJournal()); err != nil {
		t.Fatalf("SaveJournal() error = %v", err)
	}

	linkedParent := filepath.Join(root, "linked-parent")
	if err := os.Symlink(outside, linkedParent); err != nil {
		t.Fatalf("create parent directory symlink: %v", err)
	}

	_, err := LoadJournal(filepath.Join(linkedParent, "session.json"))
	if err == nil {
		t.Fatal("LoadJournal() error = nil, want symlinked ancestor rejection")
	}
	if !strings.Contains(err.Error(), "symbolic links") {
		t.Fatalf("LoadJournal() error = %q, want symbolic-link rejection", err)
	}
}
