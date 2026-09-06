package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunWithDependenciesPromptsForModelBeforeFirstTask(t *testing.T) {
	in := strings.NewReader("inspect Workspace\n/model arena/model-b\ninspect Workspace again\n/exit\n")
	var out bytes.Buffer
	var tasks []string
	modelCalls := 0

	listModels := func(context.Context) ([]string, error) {
		modelCalls++
		return []string{"arena/model-b", "arena/model-a"}, nil
	}
	runTask := func(_ context.Context, task string) error {
		tasks = append(tasks, task)
		return nil
	}

	if err := runWithDependencies(context.Background(), in, &out, nil, listModels, runTask); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	if modelCalls != 1 {
		t.Fatalf("model discovery calls = %d, want 1", modelCalls)
	}
	if len(tasks) != 1 || tasks[0] != "inspect Workspace again" {
		t.Fatalf("dispatched tasks = %#v, want only task after model selection", tasks)
	}
	got := out.String()
	if !strings.Contains(got, "Select a model with /model <id> before sending a task:\n") {
		t.Fatalf("missing model selection prompt: %q", got)
	}
	if !strings.Contains(got, "arena/model-a\narena/model-b\n") {
		t.Fatalf("prompt missing discovered models: %q", got)
	}
}
