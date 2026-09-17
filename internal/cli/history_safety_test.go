package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestCommandHandlerSanitizesHistoryDisplayText(t *testing.T) {
	var out bytes.Buffer
	history := func(context.Context) ([]string, error) {
		return []string{"read_script: safe\x1b[2Jhidden\u200btext"}, nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{History: history}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "history"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	const want = "read_script: safe[2Jhiddentext\n"
	if got := out.String(); got != want {
		t.Fatalf("history output = %q, want %q", got, want)
	}
}

func TestCommandHandlerBoundsHistoryDisplayLines(t *testing.T) {
	var out bytes.Buffer
	longEntry := "tool_call: " + strings.Repeat("界", 4096)
	history := func(context.Context) ([]string, error) {
		return []string{longEntry}, nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{History: history}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "history"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	line := strings.TrimSuffix(out.String(), "\n")
	if got := len([]rune(line)); got > maxHistoryDisplayRunes {
		t.Fatalf("history display line runes = %d, want <= %d", got, maxHistoryDisplayRunes)
	}
	if !strings.HasSuffix(line, "…") {
		t.Fatalf("history display line = %q, want truncation marker", line)
	}
}
