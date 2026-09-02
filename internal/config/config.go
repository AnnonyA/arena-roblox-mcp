package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Arena      ArenaConfig                `json:"arena"`
	Agent      AgentConfig                `json:"agent"`
	MCPServers map[string]MCPServerConfig `json:"mcpServers"`
}

type ArenaConfig struct {
	APIKeyEnv string   `json:"apiKeyEnv"`
	Model     string   `json:"model"`
	Fallbacks []string `json:"fallbacks"`
	Stream    bool     `json:"stream"`
}

type AgentConfig struct {
	MaxToolRounds int    `json:"maxToolRounds"`
	AutoPlaytest  bool   `json:"autoPlaytest"`
	ContextBudget string `json:"contextBudget"`
	SafeMode      bool   `json:"safeMode"`
}

type MCPServerConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func Default() Config {
	return Config{
		Arena: ArenaConfig{
			APIKeyEnv: "ARENA_API_KEY",
			Stream:    true,
		},
		Agent: AgentConfig{
			MaxToolRounds: 12,
			AutoPlaytest:  true,
			ContextBudget: "balanced",
			SafeMode:      true,
		},
		MCPServers: map[string]MCPServerConfig{},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.Arena.APIKeyEnv == "" {
		cfg.Arena.APIKeyEnv = "ARENA_API_KEY"
	}
	if cfg.Agent.MaxToolRounds <= 0 {
		cfg.Agent.MaxToolRounds = 12
	}
	if cfg.Agent.ContextBudget == "" {
		cfg.Agent.ContextBudget = "balanced"
	}
	return cfg, nil
}

func ResolveAPIKey(cfg Config, getenv func(string) string) (string, error) {
	name := cfg.Arena.APIKeyEnv
	if name == "" {
		name = "ARENA_API_KEY"
	}
	if key := getenv(name); key != "" {
		return key, nil
	}
	return "", errors.New("Arena API key is not configured")
}
