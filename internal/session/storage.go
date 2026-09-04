package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveJournal(path string, journal *Journal) error {
	data, err := json.Marshal(journal.Changes())
	if err != nil {
		return fmt.Errorf("encode session journal: %w", err)
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create session journal directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write session journal: %w", err)
	}
	return nil
}

func LoadJournal(path string) (*Journal, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read session journal: %w", err)
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
