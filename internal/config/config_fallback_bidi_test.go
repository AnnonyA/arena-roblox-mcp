package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsBidirectionalFormattingInArenaFallback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arena-rbx.json")
	data := `{"arena":{"fallbacks":["safe\u202emodel"]}}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected arena.fallbacks with bidirectional formatting to be rejected")
	}
	if !strings.Contains(err.Error(), "arena.fallbacks[0]") || !strings.Contains(err.Error(), "bidirectional") {
		t.Fatalf("unexpected error: %v", err)
	}
}
