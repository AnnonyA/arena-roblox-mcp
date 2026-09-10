package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
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
	if !utf8.Valid(data) {
		return Config{}, errors.New("config file must contain valid UTF-8")
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return Config{}, errors.New("config file must contain a JSON object")
	}
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return Config{}, err
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

func rejectDuplicateJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	return scanJSONValue(decoder)
}

func scanJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}

	switch delim {
	case '{':
		seen := make(map[string]struct{})
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return errors.New("config file contains an invalid JSON object key")
			}
			if _, exists := seen[key]; exists {
				return fmt.Errorf("config file contains duplicate JSON key %q", key)
			}
			seen[key] = struct{}{}
			if err := scanJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	case '[':
		for decoder.More() {
			if err := scanJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	default:
		return nil
	}
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
