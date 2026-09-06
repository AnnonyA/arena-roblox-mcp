package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsUndoReportsWhenNoChangesAreRecorded(t *testing.T) {
	in := strings.NewReader("/undo\n/exit\n")
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/test-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "Nothing to undo.\n") {
		t.Fatalf("undo command did not report empty journal: %q", got)
	}
}
