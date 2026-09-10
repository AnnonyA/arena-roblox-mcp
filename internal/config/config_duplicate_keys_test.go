package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRejectsDuplicateJSONKeys(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "top level",
			data: `{"arena":{"model":"first"},"arena":{"model":"second"}}`,
		},
		{
			name: "nested object",
			data: `{"agent":{"safeMode":true,"safeMode":false}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "arena-rbx.json")
			if err := os.WriteFile(path, []byte(tt.data), 0o600); err != nil {
				t.Fatalf("WriteFile: %v", err)
			}

			if _, err := Load(path); err == nil {
				t.Fatal("Load error = nil, want duplicate JSON key error")
			}
		})
	}
}
