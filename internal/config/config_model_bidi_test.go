package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsBidirectionalFormattingInArenaModel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := `{"arena":{"model":"safe\u202emodel"}}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected arena.model with bidirectional formatting to be rejected")
	}
	if !strings.Contains(err.Error(), "arena.model") || !strings.Contains(err.Error(), "bidirectional") {
		t.Fatalf("unexpected error: %v", err)
	}
}
