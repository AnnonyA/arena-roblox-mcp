package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsMCPCommandWithSurroundingWhitespace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := []byte(`{"mcpServers":{"Roblox_Studio":{"command":" cmd.exe ","args":["/c","mcp.bat"]}}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load error = nil, want surrounding-whitespace error")
	}
	if !strings.Contains(err.Error(), "command must not have surrounding whitespace") {
		t.Fatalf("Load error = %q, want surrounding-whitespace error", err)
	}
}
