package config

import (
	"path/filepath"
	"testing"
)

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

func TestDefaultIncludesRobloxStudioMCP(t *testing.T) {
	cfg := Default()
	server, ok := cfg.MCPServers["Roblox_Studio"]
	if !ok {
		t.Fatal("Roblox_Studio default MCP server missing")
	}
	if server.Command != "cmd.exe" {
		t.Fatalf("Command = %q", server.Command)
	}
	want := []string{"/c", `%LOCALAPPDATA%\Roblox\mcp.bat`}
	if len(server.Args) != len(want) {
		t.Fatalf("Args = %#v", server.Args)
	}
	for i := range want {
		if server.Args[i] != want[i] {
			t.Fatalf("Args[%d] = %q, want %q", i, server.Args[i], want[i])
		}
	}
}

func TestLoadMissingConfigUsesDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "arena-rbx.json"))
	if err != nil {
		t.Fatalf("Load missing config: %v", err)
	}
	if cfg.Arena.APIKeyEnv != "ARENA_API_KEY" {
		t.Fatalf("APIKeyEnv = %q", cfg.Arena.APIKeyEnv)
	}
	if _, ok := cfg.MCPServers["Roblox_Studio"]; !ok {
		t.Fatal("Roblox_Studio default MCP server missing")
	}
}
