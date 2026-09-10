package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvDoesNotApplyPartialStateOnParseError(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("FIRST=value\nBROKEN_LINE\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	setCalls := 0
	err := LoadDotEnv(
		path,
		func(string) (string, bool) { return "", false },
		func(string, string) error {
			setCalls++
			return nil
		},
	)
	if err == nil {
		t.Fatal("LoadDotEnv error = nil, want parse error")
	}
	if setCalls != 0 {
		t.Fatalf("set calls = %d, want 0 after parse error", setCalls)
	}
}
