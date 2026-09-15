package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/roblox"
)

func TestStudioListSanitizesUntrustedDisplayText(t *testing.T) {
	in := strings.NewReader("/studio\n/exit\n")
	var out bytes.Buffer
	listStudios := func(context.Context) ([]roblox.StudioSession, error) {
		return []roblox.StudioSession{
			{ID: "studio-a\x1b[31m", Name: "unsafe\u200bname", PlaceID: "123\x1b[2J"},
			{ID: "studio-b", Name: "safe", PlaceID: "456"},
		}, nil
	}

	if err := runWithStudioDependencies(context.Background(), in, &out, nil, nil, nil, nil, listStudios); err != nil {
		t.Fatalf("runWithStudioDependencies() error = %v", err)
	}
	got := out.String()
	if strings.ContainsRune(got, '\x1b') {
		t.Fatalf("terminal control character from Studio metadata reached output: %q", got)
	}
	if strings.ContainsRune(got, '\u200b') {
		t.Fatalf("Unicode format character from Studio metadata reached output: %q", got)
	}
}
