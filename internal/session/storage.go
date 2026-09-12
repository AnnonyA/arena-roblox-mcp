package session

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const maxSessionJournalBytes = 16 * 1024 * 1024

func SaveJournal(path string, journal *Journal) error {
	if journal == nil {
		return fmt.Errorf("encode session journal: nil session journal")
	}

	data, err := json.Marshal(journal.Changes())
	if err != nil {
		return fmt.Errorf("encode session journal: %w", err)
	}
	data = append(data, '\n')
	if len(data) > maxSessionJournalBytes {
		return fmt.Errorf("encode session journal: session journal exceeds %d bytes", maxSessionJournalBytes)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create session journal directory: %w", err)
	}
	if dir != "." {
		info, err := os.Lstat(dir)
		if err != nil {
			return fmt.Errorf("inspect session journal directory: %w", err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("secure session journal directory: symbolic links are not allowed")
		}
		if err := os.Chmod(dir, 0o700); err != nil {
			return fmt.Errorf("secure session journal directory permissions: %w", err)
		}
	}

	tmp, err := os.CreateTemp(dir, ".arena-rbx-journal-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary session journal: %w", err)
	}
	tmpPath := tmp.Name()
	keepTemp := true
	defer func() {
		if keepTemp {
			_ = os.Remove(tmpPath)
		}
	}()

	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("secure temporary session journal permissions: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temporary session journal: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temporary session journal: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary session journal: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace session journal: %w", err)
	}
	keepTemp = false

	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("secure session journal permissions: %w", err)
	}
	return nil
}

func LoadJournal(path string) (*Journal, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("read session journal: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("read session journal: symbolic links are not allowed")
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("read session journal: only regular files are allowed")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read session journal: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxSessionJournalBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read session journal: %w", err)
	}
	if len(data) > maxSessionJournalBytes {
		return nil, fmt.Errorf("read session journal: session journal exceeds %d bytes", maxSessionJournalBytes)
	}

	var changes []Change
	if err := json.Unmarshal(data, &changes); err != nil {
		return nil, fmt.Errorf("decode session journal: %w", err)
	}

	journal := NewJournal()
	for _, change := range changes {
		journal.Record(change)
	}
	return journal, nil
}
