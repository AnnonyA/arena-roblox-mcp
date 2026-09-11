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

func TestSaveJournalTightensExistingDirectoryPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows directory mode permissions are not POSIX")
	}

	root := t.TempDir()
	dir := filepath.Join(root, "sessions")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("seed session directory: %v", err)
	}
	path := filepath.Join(dir, "session.json")

	journal := NewJournal()
	journal.Record(Change{Tool: "edit_script", Resource: "game.ServerScriptService.Main", Reversible: true})

	if err := SaveJournal(path, journal); err != nil {
		t.Fatalf("SaveJournal() error = %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat session directory: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("session directory permissions = %o, want 700", got)
	}
}

func TestSaveJournalReplacesSymlinkWithoutFollowingIt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows symlink creation may require elevated privileges")
	}

	root := t.TempDir()
	dir := filepath.Join(root, "sessions")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatalf("create session directory: %v", err)
	}

	target := filepath.Join(root, "outside.txt")
	const original = "do not overwrite\n"
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatalf("seed symlink target: %v", err)
	}

	path := filepath.Join(dir, "session.json")
	if err := os.Symlink(target, path); err != nil {
		t.Fatalf("create journal symlink: %v", err)
	}

	journal := NewJournal()
	journal.Record(Change{Tool: "edit_script", Resource: "game.ServerScriptService.Main", Reversible: true})
	if err := SaveJournal(path, journal); err != nil {
		t.Fatalf("SaveJournal() error = %v", err)
	}

	outside, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read symlink target: %v", err)
	}
	if got := string(outside); got != original {
		t.Fatalf("symlink target content = %q, want %q", got, original)
	}

	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("lstat journal path: %v", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("journal path is still a symlink")
	}

	loaded, err := LoadJournal(path)
	if err != nil {
		t.Fatalf("LoadJournal() error = %v", err)
	}
	if got := len(loaded.Changes()); got != 1 {
		t.Fatalf("loaded change count = %d, want 1", got)
	}
}
