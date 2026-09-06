package main

import (
	"bytes"
	"context"
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
