package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsLineSeparatorsInMCPProcessConfig(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "command",
			data: `{"mcpServers":{"Roblox_Studio":{"command":"cmd\u2028.exe","args":[]}}}`,
			want: "mcpServers.Roblox_Studio.command",
		},
		{
			name: "argument",
			data: `{"mcpServers":{"Roblox_Studio":{"command":"cmd.exe","args":["/c","echo\u2029ok"]}}}`,
			want: "mcpServers.Roblox_Studio.args[1]",
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
				t.Fatal("Load error = nil, want line-separator validation error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Load error = %q, want validation path %q", err, tt.want)
			}
			if !strings.Contains(err.Error(), "line separators") {
				t.Fatalf("Load error = %q, want line-separator validation error", err)
			}
		})
	}
}
