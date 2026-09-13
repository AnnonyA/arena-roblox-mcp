package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsControlCharactersInMCPServerName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := `{"mcpServers":{"Roblox\u000aStudio":{"command":"cmd.exe","args":[]}}}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load error = nil, want MCP server-name control-character validation error")
	}
	if !strings.Contains(err.Error(), "mcpServers server name") {
		t.Fatalf("Load error = %q, want MCP server-name validation error", err)
	}
}
