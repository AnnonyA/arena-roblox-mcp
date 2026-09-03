package mcp

import "testing"

func TestSnapshotCommandConfigCopiesMutableInput(t *testing.T) {
	config := CommandConfig{
		Command:       "  roblox-mcp  ",
		Args:          []string{"--stdio"},
		ClientName:    "arena-rbx",
		ClientVersion: "0.1.0",
	}

	snapshot := snapshotCommandConfig(config)
	config.Args[0] = "--mutated"

	if snapshot.Command != "roblox-mcp" {
		t.Fatalf("snapshot.Command = %q, want %q", snapshot.Command, "roblox-mcp")
	}
	if got := snapshot.Args[0]; got != "--stdio" {
		t.Fatalf("snapshot.Args[0] = %q after caller mutation, want %q", got, "--stdio")
	}
}
