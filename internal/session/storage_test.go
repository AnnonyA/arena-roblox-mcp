package session

import (
	"path/filepath"
	"testing"
)

func TestJournalSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "session.json")
	want := Change{
		Tool:       "edit_script",
		Resource:   "game.ServerScriptService.Main",
		Before:     "print(\"before\")\n",
		After:      "print(\"after\")\n",
		Reversible: true,
	}
	journal := NewJournal()
	journal.Record(want)

	if err := SaveJournal(path, journal); err != nil {
		t.Fatalf("SaveJournal() error = %v", err)
	}

	loaded, err := LoadJournal(path)
	if err != nil {
		t.Fatalf("LoadJournal() error = %v", err)
	}
	changes := loaded.Changes()
	if len(changes) != 1 || changes[0] != want {
		t.Fatalf("loaded Changes() = %#v, want [%#v]", changes, want)
	}
}
