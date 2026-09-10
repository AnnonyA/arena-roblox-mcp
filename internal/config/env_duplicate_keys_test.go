package config

import (
	"strings"
	"testing"
)

func TestLoadDotEnvRejectsDuplicateKeysBeforeApplyingVariables(t *testing.T) {
	path := writeDotEnvFile(t, "ARENA_API_KEY=first\nARENA_API_KEY=second\n")

	setCalls := 0
	err := LoadDotEnv(path, func(string) (string, bool) {
		return "", false
	}, func(string, string) error {
		setCalls++
		return nil
	})

	if err == nil {
		t.Fatal("LoadDotEnv error = nil, want duplicate key error")
	}
	if !strings.Contains(err.Error(), "duplicate") || !strings.Contains(err.Error(), "ARENA_API_KEY") {
		t.Fatalf("LoadDotEnv error = %q, want duplicate ARENA_API_KEY error", err)
	}
	if setCalls != 0 {
		t.Fatalf("set calls = %d, want 0", setCalls)
	}
}
