package session

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSaveJournalTightensExistingFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows file mode permissions are not POSIX")
	}

	path := filepath.Join(t.TempDir(), "session.json")
	if err := os.WriteFile(path, []byte("[]\n"), 0o644); err != nil {
		t.Fatalf("seed journal: %v", err)
	}

	journal := NewJournal()
	journal.Record(Change{Tool: "edit_script", Resource: "game.ServerScriptService.Main", Reversible: true})

	if err := SaveJournal(path, journal); err != nil {
		t.Fatalf("SaveJournal() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat journal: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("journal permissions = %o, want 600", got)
	}
}
