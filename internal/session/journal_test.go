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

func TestJournalCommitUndoConsumesLatestReversibleChange(t *testing.T) {
	journal := NewJournal()
	first := Change{Tool: "edit_script", Resource: "first", Reversible: true}
	second := Change{Tool: "edit_script", Resource: "second", Reversible: true}
	journal.Record(first)
	journal.Record(second)

	if err := journal.CommitUndo(); err != nil {
		t.Fatalf("CommitUndo() error = %v", err)
	}

	changes := journal.Changes()
	if len(changes) != 1 || changes[0] != first {
		t.Fatalf("Changes() after CommitUndo() = %#v, want [%#v]", changes, first)
	}

	got, err := journal.UndoCandidate()
	if err != nil {
		t.Fatalf("UndoCandidate() after CommitUndo() error = %v", err)
	}
	if got != first {
		t.Fatalf("UndoCandidate() after CommitUndo() = %#v, want %#v", got, first)
	}
}

func TestJournalCommitUndoRejectsLatestIrreversibleChange(t *testing.T) {
	journal := NewJournal()
	change := Change{Tool: "delete_instance", Resource: "irreversible", Reversible: false}
	journal.Record(change)

	err := journal.CommitUndo()
	if !errors.Is(err, ErrLatestChangeIrreversible) {
		t.Fatalf("CommitUndo() error = %v, want ErrLatestChangeIrreversible", err)
	}
	changes := journal.Changes()
	if len(changes) != 1 || changes[0] != change {
		t.Fatalf("Changes() after rejected CommitUndo() = %#v, want [%#v]", changes, change)
	}
}
