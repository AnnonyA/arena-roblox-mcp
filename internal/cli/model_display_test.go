package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestCommandHandlerOmitsOversizedModelIDs(t *testing.T) {
	var out bytes.Buffer
	oversized := strings.Repeat("m", maxModelIDDisplayRunes+1)
	models := func(context.Context) ([]string, error) {
		return []string{"arena-code", oversized, "arena-fast"}, nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Models: models}, nil)

	if _, err := handler(context.Background(), Input{Kind: InputCommand, Command: "models"}); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if got, want := out.String(), "arena-code\narena-fast\n"; got != want {
		t.Fatalf("models output = %q, want %q", got, want)
	}
}
