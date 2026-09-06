package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestRunWithToolDependenciesStatusReflectsSuccessfulMCPConnection(t *testing.T) {
	in := strings.NewReader("/tools\n/status\n/exit\n")
	var out bytes.Buffer

	listTools := func(context.Context) ([]string, error) {
		return []string{"read_script"}, nil
	}
	runTask := func(context.Context, []arena.Message) (string, error) { return "", nil }
	if err := runWithToolDependencies(context.Background(), in, &out, []string{"--model", "arena/test-model"}, nil, listTools, runTask); err != nil {
		t.Fatalf("runWithToolDependencies() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "MCP        connected\n") {
		t.Fatalf("status command did not reflect successful MCP connection: %q", got)
	}
}
