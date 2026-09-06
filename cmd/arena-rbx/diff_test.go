package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsDiffReportsWhenNoChangesAreRecorded(t *testing.T) {
	in := strings.NewReader("/diff\n/exit\n")
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/test-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "No session changes recorded.\n") {
		t.Fatalf("diff command did not report empty journal: %q", got)
	}
}
