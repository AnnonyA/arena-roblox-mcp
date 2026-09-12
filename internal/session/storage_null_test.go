package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadJournalRejectsNullJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.json")
	if err := os.WriteFile(path, []byte("null\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := LoadJournal(path)
	if err == nil {
		t.Fatal("LoadJournal() error = nil, want null journal rejection")
	}
	if !strings.Contains(err.Error(), "session journal must be a JSON array") {
		t.Fatalf("LoadJournal() error = %q, want JSON array rejection", err)
	}
}
