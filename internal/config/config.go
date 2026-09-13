package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxConfigFileBytes      = 1 << 20
	maxConfigJSONDepth      = 64
	maxConfiguredToolRounds = 128
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
	cfg.Arena.APIKeyEnv = strings.TrimSpace(cfg.Arena.APIKeyEnv)
	if cfg.Arena.APIKeyEnv == "" {
		cfg.Arena.APIKeyEnv = "ARENA_API_KEY"
	}
	for _, r := range cfg.Arena.APIKeyEnv {
		if unicode.IsControl(r) {
			return Config{}, errors.New("arena.apiKeyEnv must not contain control characters")
		}
	}
	if strings.ContainsRune(cfg.Arena.APIKeyEnv, '=') {
		return Config{}, errors.New("arena.apiKeyEnv must not contain '='")
	}
	cfg.Arena.Model = strings.TrimSpace(cfg.Arena.Model)
	for _, r := range cfg.Arena.Model {
		if unicode.IsControl(r) {
			return Config{}, errors.New("arena.model must not contain control characters")
		}
	}
	for i, fallback := range cfg.Arena.Fallbacks {
		fallback = strings.TrimSpace(fallback)
		if fallback == "" {
			return Config{}, fmt.Errorf("arena.fallbacks[%d] must not be empty", i)
		}
		for _, r := range fallback {
			if unicode.IsControl(r) {
				return Config{}, fmt.Errorf("arena.fallbacks[%d] must not contain control characters", i)
			}
		}
		cfg.Arena.Fallbacks[i] = fallback
	}
	if cfg.Agent.MaxToolRounds <= 0 {
		cfg.Agent.MaxToolRounds = 12
	}
	if cfg.Agent.MaxToolRounds > maxConfiguredToolRounds {
		return Config{}, fmt.Errorf("agent.maxToolRounds must not exceed %d", maxConfiguredToolRounds)
	}
	for name, server := range cfg.MCPServers {
		server.Command = strings.TrimSpace(server.Command)
		if server.Command == "" {
			return Config{}, fmt.Errorf("mcpServers.%s.command must not be empty", name)
		}
		if strings.ContainsRune(server.Command, '\x00') {
			return Config{}, fmt.Errorf("mcpServers.%s.command must not contain NUL", name)
		}
		for i, arg := range server.Args {
			if strings.ContainsRune(arg, '\x00') {
				return Config{}, fmt.Errorf("mcpServers.%s.args[%d] must not contain NUL", name, i)
			}
		}
		cfg.MCPServers[name] = server
	}
	if cfg.Agent.ContextBudget == "" {
		cfg.Agent.ContextBudget = "balanced"
	}
	return cfg, nil
}

func rejectDuplicateJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	return scanJSONValue(decoder, 0)
}

func scanJSONValue(decoder *json.Decoder, depth int) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}

	depth++
	if depth > maxConfigJSONDepth {
		return errors.New("config file JSON is nested too deeply")
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
			if err := scanJSONValue(decoder, depth); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	case '[':
		for decoder.More() {
			if err := scanJSONValue(decoder, depth); err != nil {
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
