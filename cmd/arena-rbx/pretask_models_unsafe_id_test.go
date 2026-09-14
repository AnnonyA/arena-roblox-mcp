package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunFiltersUnsafeDiscoveredModelIDsBeforeFirstTask(t *testing.T) {
	in := strings.NewReader("inspect Workspace scripts\n/exit\n")
	var out bytes.Buffer

	listModels := func(context.Context) ([]string, error) {
		return []string{
			"arena/safe-model",
			"arena/escape\x1b[31m",
			"arena/zero\u200bwidth",
		}, nil
	}
	if err := runWithArgsAndModels(context.Background(), in, &out, nil, listModels); err != nil {
		t.Fatalf("runWithArgsAndModels() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "arena/safe-model") {
		t.Fatalf("safe model missing from pre-task output: %q", got)
	}
	if strings.ContainsRune(got, '\x1b') {
		t.Fatalf("control character from discovered model reached pre-task output: %q", got)
	}
	if strings.ContainsRune(got, '\u200b') {
		t.Fatalf("Unicode format character from discovered model reached pre-task output: %q", got)
	}
}
