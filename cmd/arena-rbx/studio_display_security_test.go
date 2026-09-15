package main

import (
	"errors"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/roblox"
)

func TestStudioDiscoveryRejectsUnsafeSessionID(t *testing.T) {
	_, err := roblox.ParseStudioSessions([]byte(`{"studios":[{"studio_id":"studio-\u001b[31m","name":"safe","place_id":123}]}`))
	if !errors.Is(err, roblox.ErrInvalidStudioSession) {
		t.Fatalf("ParseStudioSessions() error = %v, want ErrInvalidStudioSession", err)
	}
}

func TestStudioDiscoverySanitizesDisplayMetadata(t *testing.T) {
	sessions, err := roblox.ParseStudioSessions([]byte(`{"studios":[{"studio_id":"studio-a","name":"unsafe\u200bname\u001b[31m","place_id":123}]}`))
	if err != nil {
		t.Fatalf("ParseStudioSessions() error = %v", err)
	}
	if got, want := sessions[0].Name, "unsafename[31m"; got != want {
		t.Fatalf("Name = %q, want %q", got, want)
	}
}
