package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsRejectsControlCharactersInModelFlag(t *testing.T) {
	in := strings.NewReader("/exit\n")
	var out bytes.Buffer

	err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/\x1b[31mspoof"})
	if err == nil {
		t.Fatal("runWithArgs() error = nil, want model validation error")
	}
	if !strings.Contains(err.Error(), "control characters") {
		t.Fatalf("runWithArgs() error = %q, want control-character validation", err)
	}
	if out.Len() != 0 {
		t.Fatalf("invalid model flag wrote terminal output: %q", out.String())
	}
}
