package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestCommandHandlerBoundsAndSanitizesStatusFields(t *testing.T) {
	var out bytes.Buffer
	longModel := strings.Repeat("界", 1024) + "\x1b[31m"
	status := func() StartupStatus {
		return StartupStatus{
			Arena:   "connected\x1b[31m",
			Studio:  "offline\rspoofed",
			Model:   longModel,
			Session: "default\u202espoofed",
		}
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Status: status}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "status"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	got := out.String()
	if strings.ContainsAny(got, "\x1b\r") || strings.ContainsRune(got, '\u202e') {
		t.Fatalf("status output contains terminal control characters: %q", got)
	}
	for _, line := range strings.Split(strings.TrimSuffix(got, "\n"), "\n") {
		if len([]rune(line)) > 280 {
			t.Fatalf("status line is unbounded: %d runes", len([]rune(line)))
		}
	}
	if !strings.Contains(got, "…") {
		t.Fatalf("status output did not truncate oversized field: %q", got)
	}
}
