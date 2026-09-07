package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/roblox"
)

func TestRunWithStudioDependenciesKeepsCLIAliveWhenNoStudioIsAvailable(t *testing.T) {
	in := strings.NewReader("/studio\n/status\n/exit\n")
	var out bytes.Buffer

	listStudios := func(context.Context) ([]roblox.StudioSession, error) {
		return nil, nil
	}

	if err := runWithStudioDependencies(context.Background(), in, &out, []string{"--model", "arena/test-model"}, nil, nil, nil, listStudios); err != nil {
		t.Fatalf("runWithStudioDependencies() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "No Roblox Studio MCP session detected.\n") {
		t.Fatalf("missing unavailable Studio message: %q", got)
	}
	if !strings.Contains(got, "MCP        connected") || !strings.Contains(got, "Studio     not connected") {
		t.Fatalf("status command did not run after unavailable Studio: %q", got)
	}
}
