package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDotEnvDoesNotOverrideExistingVariable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("ARENA_API_KEY=file-key\nOTHER=value\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	existing := map[string]string{"ARENA_API_KEY": "process-key"}
	set := map[string]string{}
	lookup := func(key string) (string, bool) {
		v, ok := existing[key]
		return v, ok
	}
	setter := func(key, value string) error {
		set[key] = value
		return nil
	}

	if err := LoadDotEnv(path, lookup, setter); err != nil {
		t.Fatal(err)
	}
	if _, ok := set["ARENA_API_KEY"]; ok {
		t.Fatal("existing process variable was overwritten")
	}
	if got := set["OTHER"]; got != "value" {
		t.Fatalf("OTHER = %q", got)
	}
}

func TestLoadDotEnvRejectsUnbalancedQuotes(t *testing.T) {
	tests := []struct {
		name  string
		entry string
	}{
		{name: "double quote", entry: "ARENA_API_KEY=\"unterminated\n"},
		{name: "single quote", entry: "ARENA_API_KEY='unterminated\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, ".env")
			if err := os.WriteFile(path, []byte(tt.entry), 0o600); err != nil {
				t.Fatal(err)
			}

			setCalled := false
			err := LoadDotEnv(path, func(string) (string, bool) { return "", false }, func(string, string) error {
				setCalled = true
				return nil
			})
			if err == nil {
				t.Fatal("LoadDotEnv error = nil, want unbalanced quote error")
			}
			if !strings.Contains(err.Error(), "line 1") {
				t.Fatalf("LoadDotEnv error = %q, want line number", err)
			}
			if setCalled {
				t.Fatal("setter called for malformed .env entry")
			}
		})
	}
}
