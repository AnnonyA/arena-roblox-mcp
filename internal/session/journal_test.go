package session

import "testing"

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
