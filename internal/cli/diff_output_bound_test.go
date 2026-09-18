package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestDiffOutputIsBounded(t *testing.T) {
	var out bytes.Buffer
	handler := NewCommandHandlerWithActions(&out, CommandActions{
		Diff: func(context.Context) (string, error) {
			return strings.Repeat("界", maxDiffDisplayRunes+1024), nil
		},
	}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "diff"})
	if err != nil {
		t.Fatalf("handler() error = %v", err)
	}

	got := []rune(out.String())
	if len(got) != maxDiffDisplayRunes {
		t.Fatalf("diff output rune count = %d, want %d", len(got), maxDiffDisplayRunes)
	}
	if got[len(got)-1] != '…' {
		t.Fatalf("diff output does not end with truncation marker: %q", out.String())
	}
}
