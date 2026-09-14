package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/roblox"
)

func TestStudioCommandFiltersUnsafeDiscoveredSessionMetadata(t *testing.T) {
	in := strings.NewReader("/studio\n/status\n/exit\n")
	var out bytes.Buffer

	listStudios := func(context.Context) ([]roblox.StudioSession, error) {
		return []roblox.StudioSession{
			{ID: "studio-safe", Name: "Safe", PlaceID: "100"},
			{ID: "studio-escape\x1b[31m", Name: "Unsafe", PlaceID: "200"},
			{ID: "studio-zero\u200bwidth", Name: "Unsafe", PlaceID: "300"},
			{ID: "studio-name", Name: "Unsafe\x1b[31m", PlaceID: "400"},
			{ID: "studio-place", Name: "Unsafe", PlaceID: "500\u200b"},
		}, nil
	}

	if err := runWithStudioDependencies(context.Background(), in, &out, nil, nil, nil, nil, listStudios); err != nil {
		t.Fatalf("runWithStudioDependencies() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Studio     studio-safe\n") {
		t.Fatalf("safe Studio session was not selected: %q", got)
	}
	if strings.ContainsRune(got, '\x1b') {
		t.Fatalf("control character from discovered Studio session reached output: %q", got)
	}
	if strings.ContainsRune(got, '\u200b') {
		t.Fatalf("Unicode format character from discovered Studio session reached output: %q", got)
	}
}
