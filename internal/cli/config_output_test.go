package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCommandHandlerBoundsConfigOutput(t *testing.T) {
	var out bytes.Buffer
	config := func(context.Context) (string, error) {
		return strings.Repeat("界", maxConfigDisplayRunes+128), nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Config: config}, nil)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "config"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("config unexpectedly requested exit")
	}
	if got := utf8.RuneCountInString(out.String()); got != maxConfigDisplayRunes {
		t.Fatalf("config output runes = %d, want %d", got, maxConfigDisplayRunes)
	}
	if !strings.HasSuffix(out.String(), "…") {
		t.Fatalf("config output should end with truncation marker: %q", out.String())
	}
}
