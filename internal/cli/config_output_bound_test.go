package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestConfigOutputIsBoundedAndTerminalSafe(t *testing.T) {
	var out bytes.Buffer
	config := strings.Repeat("界", maxConfigDisplayRunes+128) + "\x1b[31mhidden\rtext\u202e"
	handler := NewCommandHandlerWithActions(&out, CommandActions{
		Config: func(context.Context) (string, error) { return config, nil },
	}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "config"})
	if err != nil {
		t.Fatalf("handler() error = %v", err)
	}
	got := out.String()
	if strings.ContainsRune(got, '\x1b') || strings.ContainsRune(got, '\r') || strings.ContainsRune(got, '\u202e') {
		t.Fatalf("config output contains terminal control/format characters: %q", got)
	}
	if len([]rune(got)) > maxConfigDisplayRunes {
		t.Fatalf("config output runes = %d, want <= %d", len([]rune(got)), maxConfigDisplayRunes)
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("config output missing truncation marker")
	}
}
