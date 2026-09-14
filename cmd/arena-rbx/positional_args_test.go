package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsRejectsUnexpectedPositionalArguments(t *testing.T) {
	var out bytes.Buffer

	err := runWithArgs(context.Background(), strings.NewReader("/exit\n"), &out, []string{"unexpected"})
	if err == nil {
		t.Fatal("runWithArgs() error = nil, want unexpected positional argument error")
	}
	if !strings.Contains(err.Error(), "unexpected positional argument") {
		t.Fatalf("runWithArgs() error = %q, want unexpected positional argument error", err)
	}
}
