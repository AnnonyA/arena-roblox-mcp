package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsUnicodeFormatCharactersInMCPConfig(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr string
	}{
		{
			name:    "server name",
			data:    `{"mcpServers":{"Roblox\u200bStudio":{"command":"cmd.exe"}}}`,
			wantErr: "mcpServers server name must not contain Unicode format characters",
		},
		{
			name:    "command",
			data:    `{"mcpServers":{"Roblox_Studio":{"command":"cmd\u200b.exe"}}}`,
			wantErr: "mcpServers.Roblox_Studio.command must not contain Unicode format characters",
		},
		{
			name:    "argument",
			data:    `{"mcpServers":{"Roblox_Studio":{"command":"cmd.exe","args":["/\u200bc"]}}}`,
			wantErr: "mcpServers.Roblox_Studio.args[0] must not contain Unicode format characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "arena-rbx.json")
			if err := os.WriteFile(path, []byte(tt.data), 0o600); err != nil {
				t.Fatalf("WriteFile: %v", err)
			}

			_, err := Load(path)
			if err == nil {
				t.Fatal("Load succeeded, want Unicode-format validation error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Load error = %q, want %q", err, tt.wantErr)
			}
		})
	}
}
