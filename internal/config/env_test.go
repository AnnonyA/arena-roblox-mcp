package config

import (
	"os"
	"path/filepath"
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
