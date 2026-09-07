package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithArgsAndModelsStatusShowsArenaConnectedAfterModelDiscovery(t *testing.T) {
	in := strings.NewReader("/models\n/status\n/exit\n")
	var out bytes.Buffer

	listModels := func(context.Context) ([]string, error) {
		return []string{"arena/test-model"}, nil
	}
	if err := runWithArgsAndModels(context.Background(), in, &out, nil, listModels); err != nil {
		t.Fatalf("runWithArgsAndModels() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Arena      connected\n") {
		t.Fatalf("status did not reflect successful Arena model discovery: %q", got)
	}
}
