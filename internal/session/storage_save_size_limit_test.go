package session

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveJournalRejectsOversizedJournal(t *testing.T) {
	journal := NewJournal()
	journal.Record(Change{
		Tool:       "script_edit",
		Resource:   "Workspace.Script",
		After:      strings.Repeat("x", maxSessionJournalBytes),
		Reversible: true,
	})

	path := filepath.Join(t.TempDir(), "journal.json")
	err := SaveJournal(path, journal)
	if err == nil {
		t.Fatal("SaveJournal() error = nil, want oversized journal error")
	}
	if !strings.Contains(err.Error(), "session journal exceeds") {
		t.Fatalf("SaveJournal() error = %q, want size-limit error", err)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("journal file exists after rejected save; stat error = %v", statErr)
	}
}
