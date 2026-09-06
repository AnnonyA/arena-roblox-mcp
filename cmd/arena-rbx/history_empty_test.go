package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsHistoryReportsWhenNoActionsAreRecorded(t *testing.T) {
	in := strings.NewReader("/history\n/exit\n")
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/test-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "No session history recorded.\n") {
		t.Fatalf("history command did not report empty session history: %q", got)
	}
}
