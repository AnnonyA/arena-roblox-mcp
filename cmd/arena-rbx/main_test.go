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
