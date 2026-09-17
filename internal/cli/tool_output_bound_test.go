package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestToolsCommandBoundsIndividualDisplayLines(t *testing.T) {
	var out bytes.Buffer
	longDescription := strings.Repeat("x", 4096)
	handler := NewCommandHandlerWithActions(&out, CommandActions{
		Tools: func(context.Context) ([]string, error) {
			return []string{"tool — " + longDescription}, nil
		},
	}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "tools"})
	if err != nil {
		t.Fatalf("handler error = %v", err)
	}
	line := strings.TrimSuffix(out.String(), "\n")
	if got := len([]rune(line)); got > maxToolDisplayRunes {
		t.Fatalf("tool display line runes = %d, want <= %d", got, maxToolDisplayRunes)
	}
	if !strings.HasSuffix(line, "…") {
		t.Fatalf("truncated tool display = %q, want ellipsis suffix", line)
	}
}
