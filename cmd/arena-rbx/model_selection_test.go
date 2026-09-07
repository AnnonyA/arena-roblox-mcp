package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
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
	runTask := func(_ context.Context, messages []arena.Message) (string, error) {
		if len(messages) != 1 {
			t.Fatalf("messages = %#v, want one current task", messages)
		}
		tasks = append(tasks, messages[0].Content)
		return "", nil
	}

	if err := runWithDependencies(context.Background(), in, &out, nil, listModels, runTask); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	if modelCalls != 2 {
		t.Fatalf("model discovery calls = %d, want 2", modelCalls)
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

func TestRunWithDependenciesRejectsUnknownModelAndKeepsCLIUsable(t *testing.T) {
	in := strings.NewReader("/model arena/missing\ninspect Workspace\n/model arena/model-a\ninspect Workspace again\n/exit\n")
	var out bytes.Buffer
	var tasks []string

	listModels := func(context.Context) ([]string, error) {
		return []string{"arena/model-b", "arena/model-a"}, nil
	}
	runTask := func(_ context.Context, messages []arena.Message) (string, error) {
		tasks = append(tasks, messages[len(messages)-1].Content)
		return "", nil
	}

	if err := runWithDependencies(context.Background(), in, &out, nil, listModels, runTask); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	if len(tasks) != 1 || tasks[0] != "inspect Workspace again" {
		t.Fatalf("dispatched tasks = %#v, want only task after valid model selection", tasks)
	}
	if got := out.String(); !strings.Contains(got, "Error: Arena model not found: arena/missing\n") {
		t.Fatalf("missing unknown-model error: %q", got)
	}
}

func TestRunWithDependenciesModelWithoutIDListsChoicesAndKeepsCLIUsable(t *testing.T) {
	in := strings.NewReader("/model\n/status\n/exit\n")
	var out bytes.Buffer

	listModels := func(context.Context) ([]string, error) {
		return []string{"arena/model-b", "arena/model-a"}, nil
	}

	if err := runWithDependencies(context.Background(), in, &out, nil, listModels, nil); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	got := out.String()
	if strings.Contains(got, "Error: model ID is required") {
		t.Fatalf("/model without ID should guide selection, got: %q", got)
	}
	if !strings.Contains(got, "Select a model with /model <id>:\n") {
		t.Fatalf("missing model selection guidance: %q", got)
	}
	if !strings.Contains(got, "arena/model-a\narena/model-b\n") {
		t.Fatalf("missing discovered models: %q", got)
	}
	if !strings.Contains(got, "Model      not selected\n") {
		t.Fatalf("CLI did not remain usable after /model: %q", got)
	}
}
