package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsAndModelsFiltersUnsafeDiscoveredModelIDs(t *testing.T) {
	in := strings.NewReader("/models\n/exit\n")
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
		t.Fatalf("safe model missing from output: %q", got)
	}
	if strings.ContainsRune(got, '\x1b') {
		t.Fatalf("control character from discovered model reached output: %q", got)
	}
	if strings.ContainsRune(got, '\u200b') {
		t.Fatalf("Unicode format character from discovered model reached output: %q", got)
	}
}
