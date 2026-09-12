package session

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadJournalRejectsNonRegularFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "journal-dir")
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatalf("create journal directory: %v", err)
	}

	_, err := LoadJournal(path)
	if err == nil {
		t.Fatal("LoadJournal() error = nil, want non-regular file rejection")
	}
	if !strings.Contains(err.Error(), "regular file") {
		t.Fatalf("LoadJournal() error = %q, want regular-file rejection", err)
	}
}
