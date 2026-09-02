package config

import "testing"

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Arena.APIKeyEnv != "ARENA_API_KEY" {
		t.Fatalf("APIKeyEnv = %q", cfg.Arena.APIKeyEnv)
	}
	if cfg.Agent.MaxToolRounds != 12 {
		t.Fatalf("MaxToolRounds = %d", cfg.Agent.MaxToolRounds)
	}
	if !cfg.Agent.SafeMode {
		t.Fatal("SafeMode must default to true")
	}
}
