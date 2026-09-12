package session

import (
	"strings"
	"testing"
)

func TestDiffEscapesUnsafeResourceLabelCharacters(t *testing.T) {
	t.Parallel()

	journal := NewJournal()
	journal.Record(Change{
		Resource: "ServerScriptService.Main\n+++ forged\x1b[31m",
		Before:   "old",
		After:    "new",
	})

	diff := journal.Diff()
	if strings.Contains(diff, "Main\n+++ forged") {
		t.Fatalf("Diff() contains injected diff header: %q", diff)
	}
	if strings.ContainsRune(diff, '\x1b') {
		t.Fatalf("Diff() contains raw terminal escape: %q", diff)
	}
	if !strings.Contains(diff, `ServerScriptService.Main\n+++ forged\x1b[31m`) {
		t.Fatalf("Diff() = %q, want escaped resource label", diff)
	}
}
