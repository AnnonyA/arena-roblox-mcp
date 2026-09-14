package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsRejectsUnicodeFormatCharactersInModelFlag(t *testing.T) {
	in := strings.NewReader("/exit\n")
	var out bytes.Buffer

	err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/\u200bspoof"})
	if err == nil {
		t.Fatal("runWithArgs() error = nil, want model validation error")
	}
	if !strings.Contains(err.Error(), "format characters") {
		t.Fatalf("runWithArgs() error = %q, want Unicode format-character validation", err)
	}
	if out.Len() != 0 {
		t.Fatalf("invalid model flag wrote terminal output: %q", out.String())
	}
}
