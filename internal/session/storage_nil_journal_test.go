package session

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveJournalRejectsNilJournal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "journal.json")

	err := SaveJournal(path, nil)
	if err == nil {
		t.Fatal("SaveJournal() error = nil, want nil journal error")
	}
	if !strings.Contains(err.Error(), "nil session journal") {
		t.Fatalf("SaveJournal() error = %q, want nil-journal error", err)
	}
}
