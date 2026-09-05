package cli

import "testing"

func TestStartupTextShowsConnectionStateAndPrompt(t *testing.T) {
	got := StartupText(StartupStatus{
		Arena:   "connected",
		Studio:  "connected",
		Model:   "arena-code",
		Session: "default",
	})
	want := "Arena Roblox MCP\n" +
		"────────────────────────────\n" +
		"Arena      connected\n" +
		"Studio     connected\n" +
		"Model      arena-code\n" +
		"Session    default\n\n" +
		"> "
	if got != want {
		t.Fatalf("StartupText() = %q, want %q", got, want)
	}
}
