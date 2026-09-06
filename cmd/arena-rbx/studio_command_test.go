package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/roblox"
)

func TestRunWithStudioDependenciesSelectsStudioAndUpdatesStatus(t *testing.T) {
	in := strings.NewReader("/studio studio-b\n/status\n/exit\n")
	var out bytes.Buffer
	calls := 0

	listStudios := func(context.Context) ([]roblox.StudioSession, error) {
		calls++
		return []roblox.StudioSession{
			{ID: "studio-a", Name: "Studio A"},
			{ID: "studio-b", Name: "Studio B"},
		}, nil
	}

	if err := runWithStudioDependencies(context.Background(), in, &out, []string{"--model", "arena/test-model"}, nil, nil, nil, listStudios); err != nil {
		t.Fatalf("runWithStudioDependencies() error = %v", err)
	}
	if calls != 1 {
		t.Fatalf("studio discovery calls = %d, want 1", calls)
	}
	if got := out.String(); !strings.Contains(got, "Studio     studio-b\n") {
		t.Fatalf("status command did not reflect /studio selection: %q", got)
	}
}

func TestRunWithStudioDependenciesAutoSelectsOnlyStudio(t *testing.T) {
	in := strings.NewReader("/studio\n/status\n/exit\n")
	var out bytes.Buffer

	listStudios := func(context.Context) ([]roblox.StudioSession, error) {
		return []roblox.StudioSession{{ID: "only-studio", Name: "Only Studio"}}, nil
	}

	if err := runWithStudioDependencies(context.Background(), in, &out, []string{"--model", "arena/test-model"}, nil, nil, nil, listStudios); err != nil {
		t.Fatalf("runWithStudioDependencies() error = %v", err)
	}
	if got := out.String(); !strings.Contains(got, "Studio     only-studio\n") {
		t.Fatalf("status command did not auto-select the only Studio: %q", got)
	}
}
