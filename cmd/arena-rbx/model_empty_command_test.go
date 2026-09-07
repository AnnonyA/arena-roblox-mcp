package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithDependenciesModelWithoutIDReportsEmptyDiscoveryWithoutSelectionPrompt(t *testing.T) {
	in := strings.NewReader("/model\n/status\n/exit\n")
	var out bytes.Buffer

	listModels := func(context.Context) ([]string, error) {
		return []string{}, nil
	}

	if err := runWithDependencies(context.Background(), in, &out, nil, listModels, nil); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "No Arena models available.\n") {
		t.Fatalf("missing empty-model message: %q", got)
	}
	if strings.Contains(got, "Select a model with /model <id>:\n") {
		t.Fatalf("/model should not offer an impossible selection when discovery is empty: %q", got)
	}
	if !strings.Contains(got, "Model      not selected\n") {
		t.Fatalf("CLI did not remain usable after empty /model discovery: %q", got)
	}
}
