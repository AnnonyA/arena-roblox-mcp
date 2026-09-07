package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunReportsWhenNoModelsAvailableBeforeTask(t *testing.T) {
	in := strings.NewReader("inspect Workspace scripts\n/exit\n")
	var out bytes.Buffer

	listModels := func(context.Context) ([]string, error) {
		return nil, nil
	}

	if err := runWithArgsAndModels(context.Background(), in, &out, nil, listModels); err != nil {
		t.Fatalf("runWithArgsAndModels() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "No Arena models available.\n") {
		t.Fatalf("missing empty model discovery message: %q", got)
	}
}
