package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsExcessiveMaxToolRounds(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	if err := os.WriteFile(path, []byte(`{"agent":{"maxToolRounds":129}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load error = nil, want excessive maxToolRounds error")
	}
	if !strings.Contains(err.Error(), "maxToolRounds") || !strings.Contains(err.Error(), "128") {
		t.Fatalf("Load error = %q, want maxToolRounds limit of 128", err)
	}
}
