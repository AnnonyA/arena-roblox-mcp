package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsModelsCommandReportsEmptyDiscovery(t *testing.T) {
	in := strings.NewReader("/models\n/exit\n")
	var out bytes.Buffer

	listModels := func(context.Context) ([]string, error) { return nil, nil }
	if err := runWithArgsAndModels(context.Background(), in, &out, nil, listModels); err != nil {
		t.Fatalf("runWithArgsAndModels() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "No Arena models available.\n") {
		t.Fatalf("empty model discovery was silent: %q", got)
	}
}
