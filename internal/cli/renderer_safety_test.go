package cli

import (
	"strings"
	"testing"
)

func TestStatusTextSanitizesDynamicFields(t *testing.T) {
	status := StartupStatus{
		Arena:   "connected\x1b[2J",
		MCP:     "connected\u200b",
		Studio:  "Studio\x00One",
		Model:   "safe\x1b[31m",
		Session: "default\u200b",
	}

	got := StatusText(status)
	for _, unsafe := range []string{"\x1b", "\u200b", "\x00"} {
		if strings.Contains(got, unsafe) {
			t.Fatalf("status output contains unsafe terminal character %q: %q", unsafe, got)
		}
	}
	for _, want := range []string{"connected[2J", "StudioOne", "safe[31m", "default"} {
		if !strings.Contains(got, want) {
			t.Fatalf("status output %q does not contain %q", got, want)
		}
	}
}
