package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const maxConfigFileBytes = 1 << 20

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
		MCPServers: map[string]MCPServerConfig{
			"Roblox_Studio": {
				Command: "cmd.exe",
				Args: []string{
					"/c",
					`%LOCALAPPDATA%\Roblox\mcp.bat`,
				},
			},
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return Config{}, err
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxConfigFileBytes+1))
	if err != nil {
		return Config{}, err
	}
	if len(data) > maxConfigFileBytes {
		return Config{}, fmt.Errorf("config file is too large: maximum size is %d bytes", maxConfigFileBytes)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return Config{}, errors.New("config file must contain a single JSON object")
		}
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
