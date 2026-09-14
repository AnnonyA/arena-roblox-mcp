package session

import "testing"

func TestJournalCommitUndoClearsRemovedChange(t *testing.T) {
	journal := NewJournal()
	journal.Record(Change{
		Tool:       "edit_script",
		Resource:   "ServerScriptService.Main",
		Before:     "large before snapshot",
		After:      "large after snapshot",
		Reversible: true,
	})

	if err := journal.CommitUndo(); err != nil {
		t.Fatalf("CommitUndo() error = %v", err)
	}
	if len(journal.changes) != 0 {
		t.Fatalf("len(changes) = %d, want 0", len(journal.changes))
	}

	retained := journal.changes[:cap(journal.changes)]
	if len(retained) != 1 {
		t.Fatalf("backing slice len = %d, want 1", len(retained))
	}
	if retained[0] != (Change{}) {
		t.Fatalf("removed change retained in backing array: %#v", retained[0])
	}
}
