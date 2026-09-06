package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithToolDependenciesToolsCommandListsAvailableMCPTools(t *testing.T) {
	in := strings.NewReader("/tools\n/exit\n")
	var out bytes.Buffer
	calls := 0

	listTools := func(context.Context) ([]string, error) {
		calls++
		return []string{"script.read — Read a script", "script.edit — Edit a script"}, nil
	}
	if err := runWithToolDependencies(context.Background(), in, &out, nil, nil, listTools, nil); err != nil {
		t.Fatalf("runWithToolDependencies() error = %v", err)
	}

	if calls != 1 {
		t.Fatalf("tool discovery calls = %d, want 1", calls)
	}
	got := out.String()
	if !strings.Contains(got, "script.read — Read a script\nscript.edit — Edit a script\n") {
		t.Fatalf("tools command missing discovered MCP tools: %q", got)
	}
}

func TestRunWithToolDependenciesToolsCommandReportsEmptyRegistry(t *testing.T) {
	in := strings.NewReader("/tools\n/exit\n")
	var out bytes.Buffer

	listTools := func(context.Context) ([]string, error) {
		return nil, nil
	}
	if err := runWithToolDependencies(context.Background(), in, &out, nil, nil, listTools, nil); err != nil {
		t.Fatalf("runWithToolDependencies() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "No MCP tools available.\n") {
		t.Fatalf("tools command did not report empty registry: %q", got)
	}
}
