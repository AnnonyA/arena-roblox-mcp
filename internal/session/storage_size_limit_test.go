package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadJournalRejectsOversizedFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "journal.json")
	data := "[" + strings.Repeat(" ", 16*1024*1024) + "]"
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadJournal(path); err == nil {
		t.Fatal("LoadJournal error = nil, want oversized journal error")
	} else if !strings.Contains(err.Error(), "session journal exceeds") {
		t.Fatalf("LoadJournal error = %q, want size-limit error", err)
	}
}
