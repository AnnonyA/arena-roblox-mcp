package mcp

import "testing"

func TestNewCommandClientRejectsMissingCommand(t *testing.T) {
	client, err := NewCommandClient(CommandConfig{})
	if err == nil {
		t.Fatal("NewCommandClient error = nil, want missing command error")
	}
	if client != nil {
		t.Fatal("NewCommandClient returned a client for an empty command")
	}
}
