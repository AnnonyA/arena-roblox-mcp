package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestHistorySanitizesUnsafeDisplayCharactersWithoutMutatingInput(t *testing.T) {
	history := []string{"task: inspect\x1b[31m hidden\u200b text"}
	original := append([]string(nil), history...)
	var out bytes.Buffer
	handler := NewCommandHandlerWithActions(&out, CommandActions{
		History: func(context.Context) ([]string, error) { return history, nil },
	}, nil)

	if _, err := handler(context.Background(), Input{Kind: InputCommand, Command: "history"}); err != nil {
		t.Fatalf("history command error = %v", err)
	}
	got := out.String()
	if strings.ContainsRune(got, '\x1b') || strings.ContainsRune(got, '\u200b') {
		t.Fatalf("history command emitted unsafe display characters: %q", got)
	}
	if !strings.Contains(got, "task: inspect[31m hidden text\n") {
		t.Fatalf("history command did not preserve safe text: %q", got)
	}
	if history[0] != original[0] {
		t.Fatalf("history command mutated caller-owned data: got %q, want %q", history, original)
	}
}
