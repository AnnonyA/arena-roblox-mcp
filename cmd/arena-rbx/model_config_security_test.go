package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWithArgsRejectsUnsafeConfiguredModelBeforeRendering(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "arena-rbx.json"), []byte(`{"arena":{"model":"arena/unsafe\u200bmodel"}}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldwd) })

	var out bytes.Buffer
	err = runWithArgs(context.Background(), strings.NewReader("/exit\n"), &out, nil)
	if err == nil || !strings.Contains(err.Error(), "Unicode format") {
		t.Fatalf("runWithArgs() error = %v, want Unicode format rejection", err)
	}
	if got := out.String(); got != "" {
		t.Fatalf("unsafe configured model reached terminal output: %q", got)
	}
}
