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
	tests := []struct {
		name      string
		modelJSON string
		wantError string
	}{
		{name: "unicode format", modelJSON: `arena/unsafe\u200bmodel`, wantError: "Unicode format"},
		{name: "control character", modelJSON: `arena/unsafe\u001bmodel`, wantError: "control characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			configJSON := `{"arena":{"model":"` + tt.modelJSON + `"}}`
			if err := os.WriteFile(filepath.Join(dir, "arena-rbx.json"), []byte(configJSON), 0o600); err != nil {
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
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("runWithArgs() error = %v, want %q rejection", err, tt.wantError)
			}
			if got := out.String(); got != "" {
				t.Fatalf("unsafe configured model reached terminal output: %q", got)
			}
		})
	}
}
