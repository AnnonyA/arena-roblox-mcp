package mcp

import (
	"errors"
	"testing"
)

func TestNewCommandClientRejectsMissingCommand(t *testing.T) {
	client, err := NewCommandClient(CommandConfig{})
	if err == nil {
		t.Fatal("NewCommandClient error = nil, want missing command error")
	}
	if client != nil {
		t.Fatal("NewCommandClient returned a client for an empty command")
	}
}

func TestSnapshotCommandConfigNormalizesClientMetadata(t *testing.T) {
	config := snapshotCommandConfig(CommandConfig{
		Command:       "  roblox-mcp  ",
		ClientName:    "  arena-rbx  ",
		ClientVersion: "  0.1.0  ",
	})

	if config.ClientName != "arena-rbx" {
		t.Fatalf("ClientName = %q, want %q", config.ClientName, "arena-rbx")
	}
	if config.ClientVersion != "0.1.0" {
		t.Fatalf("ClientVersion = %q, want %q", config.ClientVersion, "0.1.0")
	}
}

func TestValidateToolPageCursorRejectsRepeatedCursor(t *testing.T) {
	seen := map[string]struct{}{}
	if err := validateToolPageCursor(1, seen, "next"); err != nil {
		t.Fatalf("first cursor rejected: %v", err)
	}
	if err := validateToolPageCursor(2, seen, "next"); !errors.Is(err, ErrToolPaginationCycle) {
		t.Fatalf("repeated cursor error = %v, want %v", err, ErrToolPaginationCycle)
	}
}

func TestValidateToolPageCursorRejectsExcessivePagination(t *testing.T) {
	seen := map[string]struct{}{}
	if err := validateToolPageCursor(maxToolListPages+1, seen, "still-more"); !errors.Is(err, ErrToolPaginationLimit) {
		t.Fatalf("pagination limit error = %v, want %v", err, ErrToolPaginationLimit)
	}
}
