package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestRunWithDependenciesStatusShowsArenaConnectedAfterSuccessfulTask(t *testing.T) {
	in := strings.NewReader("inspect Workspace scripts\n/status\n/exit\n")
	var out bytes.Buffer

	runTask := func(context.Context, []arena.Message) (string, error) { return "done", nil }
	if err := runWithDependencies(context.Background(), in, &out, []string{"--model", "arena/test-model"}, nil, runTask); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Arena      connected\n") {
		t.Fatalf("status did not reflect successful Arena request: %q", got)
	}
}
