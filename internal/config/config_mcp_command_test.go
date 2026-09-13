package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsEmptyMCPServerCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"mcpServers":{"Roblox_Studio":{"command":"   ","args":[]}}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load error = nil, want empty MCP command error")
	}
	if !strings.Contains(err.Error(), "mcpServers.Roblox_Studio.command") {
		t.Fatalf("Load error = %q, want MCP command validation error", err)
	}
}

func TestLoadTrimsMCPServerCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"mcpServers":{"Roblox_Studio":{"command":"  cmd.exe  ","args":["/c","echo ok"]}}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := cfg.MCPServers["Roblox_Studio"].Command; got != "cmd.exe" {
		t.Fatalf("MCP command = %q, want %q", got, "cmd.exe")
	}
}
