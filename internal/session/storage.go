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

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create session journal directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write session journal: %w", err)
	}
	return nil
}

func LoadJournal(path string) (*Journal, error) {
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
