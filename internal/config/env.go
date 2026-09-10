package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const maxDotEnvBytes = 1 << 20

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
		if len(value) > 0 && (value[0] == '"' || value[0] == '\'') {
			if len(value) < 2 || value[len(value)-1] != value[0] {
				return fmt.Errorf("invalid .env entry on line %d: unbalanced quotes", lineNo)
			}
			value = value[1 : len(value)-1]
		}
		if _, exists := lookup(key); exists {
			continue
		}
		if err := set(key, value); err != nil {
			return err
		}
	}
	return scanner.Err()
}
