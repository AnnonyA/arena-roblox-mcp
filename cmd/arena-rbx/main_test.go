package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunProvidesUsableHelpAndExitShell(t *testing.T) {
	in := strings.NewReader("/help\n/exit\n")
	var out bytes.Buffer

	if err := run(context.Background(), in, &out); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Arena Roblox MCP\n") {
		t.Fatalf("output missing startup banner: %q", got)
	}
	if !strings.Contains(got, "/help  show commands\n") {
		t.Fatalf("output missing help text: %q", got)
	}
}

func TestRunWithArgsModelFlagShowsSelectedModel(t *testing.T) {
	in := strings.NewReader("/exit\n")
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/test-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "Model      arena/test-model\n") {
		t.Fatalf("output missing selected model: %q", got)
	}
}

func TestRunWithArgsUsesConfiguredModelWhenFlagMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "arena-rbx.json"), []byte(`{"arena":{"model":"arena/config-model"}}`), 0o600); err != nil {
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

	in := strings.NewReader("/exit\n")
	var out bytes.Buffer
	if err := runWithArgs(context.Background(), in, &out, nil); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "Model      arena/config-model\n") {
		t.Fatalf("output missing configured model: %q", got)
	}
}
