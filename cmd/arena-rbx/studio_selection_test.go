package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/roblox"
)

func TestRunWithStudioDependenciesListsSessionsWhenSelectionRequired(t *testing.T) {
	in := strings.NewReader("/studio\n/exit\n")
	var out bytes.Buffer

	listStudios := func(context.Context) ([]roblox.StudioSession, error) {
		return []roblox.StudioSession{
			{ID: "studio-b", Name: "Second", PlaceID: "200"},
			{ID: "studio-a", Name: "First", PlaceID: "100"},
		}, nil
	}

	if err := runWithStudioDependencies(context.Background(), in, &out, nil, nil, nil, nil, listStudios); err != nil {
		t.Fatalf("runWithStudioDependencies() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Multiple Roblox Studio sessions detected. Select one with /studio <studio_id>:\n") {
		t.Fatalf("studio command missing selection guidance: %q", got)
	}
	if !strings.Contains(got, "studio-a — First (place 100)\nstudio-b — Second (place 200)\n") {
		t.Fatalf("studio command missing stable session list: %q", got)
	}
	if strings.Contains(got, "Error: multiple Roblox Studio sessions detected") {
		t.Fatalf("studio command surfaced generic selection error instead of choices: %q", got)
	}
}
