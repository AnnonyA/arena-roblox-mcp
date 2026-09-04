package session

import (
	"errors"
	"testing"
)

func TestJournalRecordsReversibleChange(t *testing.T) {
	journal := NewJournal()
	change := Change{
		Tool:       "edit_script",
		Resource:   "game.ServerScriptService.Main",
		Before:     "print(\"before\")",
		After:      "print(\"after\")",
		Reversible: true,
	}

	journal.Record(change)

	got := journal.Changes()
	if len(got) != 1 {
		t.Fatalf("len(Changes()) = %d, want 1", len(got))
	}
	if got[0] != change {
		t.Fatalf("Changes()[0] = %#v, want %#v", got[0], change)
	}
}

func TestJournalDiffFormatsReversibleScriptChange(t *testing.T) {
	journal := NewJournal()
	journal.Record(Change{
		Tool:       "edit_script",
		Resource:   "game.ServerScriptService.Main",
		Before:     "print(\"before\")\n",
		After:      "print(\"after\")\n",
		Reversible: true,
	})

	const want = "--- game.ServerScriptService.Main (before)\n" +
		"+++ game.ServerScriptService.Main (after)\n" +
		"@@\n" +
		"-print(\"before\")\n" +
		"+print(\"after\")\n"

	if got := journal.Diff(); got != want {
		t.Fatalf("Diff() = %q, want %q", got, want)
	}
}

func TestJournalUndoCandidateReturnsLatestChange(t *testing.T) {
	journal := NewJournal()
	journal.Record(Change{Tool: "edit_script", Resource: "first", Reversible: true})
	want := Change{Tool: "edit_script", Resource: "second", Before: "old", After: "new", Reversible: true}
	journal.Record(want)

	got, err := journal.UndoCandidate()
	if err != nil {
		t.Fatalf("UndoCandidate() error = %v", err)
	}
	if got != want {
		t.Fatalf("UndoCandidate() = %#v, want %#v", got, want)
	}
}

func TestJournalUndoCandidateRejectsLatestIrreversibleChange(t *testing.T) {
	journal := NewJournal()
	journal.Record(Change{Tool: "edit_script", Resource: "reversible", Reversible: true})
	journal.Record(Change{Tool: "delete_instance", Resource: "irreversible", Reversible: false})

	_, err := journal.UndoCandidate()
	if !errors.Is(err, ErrLatestChangeIrreversible) {
		t.Fatalf("UndoCandidate() error = %v, want ErrLatestChangeIrreversible", err)
	}
}

func TestJournalUndoCandidateRejectsEmptyJournal(t *testing.T) {
	journal := NewJournal()

	_, err := journal.UndoCandidate()
	if !errors.Is(err, ErrNoChanges) {
		t.Fatalf("UndoCandidate() error = %v, want ErrNoChanges", err)
	}
}
