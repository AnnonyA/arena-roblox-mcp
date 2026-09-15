package cli

import (
	"bytes"
	"context"
	"testing"
)

func TestCommandHandlerSanitizesDiffDisplayText(t *testing.T) {
	var out bytes.Buffer
	diff := func(context.Context) (string, error) {
		return "--- before\n+++ after\n-old\x1b[2J\n+new\u200b\n", nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Diff: diff}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "diff"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	const want = "--- before\n+++ after\n-old[2J\n+new\n"
	if got := out.String(); got != want {
		t.Fatalf("diff output = %q, want %q", got, want)
	}
}

func TestCommandHandlerSanitizesConfigDisplayText(t *testing.T) {
	var out bytes.Buffer
	config := func(context.Context) (string, error) {
		return "arena.model=safe\x1b[2J\nagent.context=balanced\u200b\n", nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Config: config}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "config"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	const want = "arena.model=safe[2J\nagent.context=balanced\n"
	if got := out.String(); got != want {
		t.Fatalf("config output = %q, want %q", got, want)
	}
}
