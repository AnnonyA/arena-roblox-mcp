package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

const maxDotEnvBytes = 1 << 20

type dotEnvEntry struct {
	key   string
	value string
}

func LoadDotEnv(path string, lookup func(string) (string, bool), set func(string, string) error) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxDotEnvBytes+1))
	if err != nil {
		return err
	}
	if len(data) > maxDotEnvBytes {
		return fmt.Errorf(".env file too large (maximum %d bytes)", maxDotEnvBytes)
	}
	if !utf8.Valid(data) {
		return fmt.Errorf(".env file contains invalid UTF-8")
	}
	if bytes.IndexByte(data, 0) >= 0 {
		return fmt.Errorf(".env file contains NUL byte")
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	var entries []dotEnvEntry
	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid .env entry on line %d", lineNo)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return fmt.Errorf("invalid .env entry on line %d", lineNo)
		}
		if _, exists := seen[key]; exists {
			return fmt.Errorf("duplicate .env key %q on line %d", key, lineNo)
		}
		seen[key] = struct{}{}
		if len(value) > 0 && (value[0] == '"' || value[0] == '\'') {
			if len(value) < 2 || value[len(value)-1] != value[0] {
				return fmt.Errorf("invalid .env entry on line %d: unbalanced quotes", lineNo)
			}
			value = value[1 : len(value)-1]
		}
		entries = append(entries, dotEnvEntry{key: key, value: value})
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	for _, entry := range entries {
		if _, exists := lookup(entry.key); exists {
			continue
		}
		if err := set(entry.key, entry.value); err != nil {
			return err
		}
	}
	return nil
}
