package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsControlCharactersInMCPProcessConfig(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "command",
			data: `{"mcpServers":{"Roblox_Studio":{"command":"cmd\u000a.exe","args":[]}}}`,
			want: "mcpServers.Roblox_Studio.command",
		},
		{
			name: "argument",
			data: `{"mcpServers":{"Roblox_Studio":{"command":"cmd.exe","args":["/c","echo\u000aok"]}}}`,
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
				t.Fatal("Load error = nil, want control-character validation error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Load error = %q, want validation path %q", err, tt.want)
			}
		})
	}
}
