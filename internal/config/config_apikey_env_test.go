package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsWhitespaceAPIKeyEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"arena":{"apiKeyEnv":"   "}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Arena.APIKeyEnv != "ARENA_API_KEY" {
		t.Fatalf("Arena.APIKeyEnv = %q, want ARENA_API_KEY", cfg.Arena.APIKeyEnv)
	}
}

func TestLoadTrimsAPIKeyEnvWhitespace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"arena":{"apiKeyEnv":"  CUSTOM_ARENA_KEY  "}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Arena.APIKeyEnv != "CUSTOM_ARENA_KEY" {
		t.Fatalf("Arena.APIKeyEnv = %q, want CUSTOM_ARENA_KEY", cfg.Arena.APIKeyEnv)
	}
}
