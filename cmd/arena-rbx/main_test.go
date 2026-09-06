package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunProvidesUsableHelpAndExitShell(t *testing.T) {
	in := strings.NewReader("/help\n/exit\n")
	var out bytes.Buffer

	if err := run(context.Background(), in, &out); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Arena Roblox MCP\n") {
		t.Fatalf("output missing startup banner: %q", got)
	}
	if !strings.Contains(got, "/help  show commands\n") {
		t.Fatalf("output missing help text: %q", got)
	}
}

func TestRunWithArgsHelpFlagPrintsHelpWithoutStartingShell(t *testing.T) {
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), strings.NewReader(""), &out, []string{"--help"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "/help  show commands\n") {
		t.Fatalf("help flag missing command help: %q", got)
	}
	if strings.Contains(got, "Arena Roblox MCP\n") {
		t.Fatalf("help flag unexpectedly started interactive shell: %q", got)
	}
}

func TestRunWithArgsModelFlagShowsSelectedModel(t *testing.T) {
	in := strings.NewReader("/exit\n")
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/test-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "Model      arena/test-model\n") {
		t.Fatalf("output missing selected model: %q", got)
	}
}

func TestRunWithArgsUsesConfiguredModelWhenFlagMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "arena-rbx.json"), []byte(`{"arena":{"model":"arena/config-model"}}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldwd) })

	in := strings.NewReader("/exit\n")
	var out bytes.Buffer
	if err := runWithArgs(context.Background(), in, &out, nil); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "Model      arena/config-model\n") {
		t.Fatalf("output missing configured model: %q", got)
	}
}

func TestRunWithArgsStatusCommandShowsCurrentStartupStatus(t *testing.T) {
	in := strings.NewReader("/status\n/exit\n")
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/status-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	got := out.String()
	if count := strings.Count(got, "Model      arena/status-model\n"); count != 2 {
		t.Fatalf("status command did not render current status; model line count = %d, output = %q", count, got)
	}
}

func TestRunWithArgsModelCommandUpdatesCurrentModel(t *testing.T) {
	in := strings.NewReader("/model arena/new-model\n/status\n/exit\n")
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/old-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Model      arena/new-model\n") {
		t.Fatalf("status command did not reflect /model selection: %q", got)
	}
}

func TestRunWithArgsModelCommandRejectsEmptyModelAndPreservesSelection(t *testing.T) {
	in := strings.NewReader("/model\n/status\n/exit\n")
	var out bytes.Buffer

	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/old-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "Error: model ID is required\n") {
		t.Fatalf("empty /model did not report actionable error: %q", got)
	}
	if count := strings.Count(got, "Model      arena/old-model\n"); count != 2 {
		t.Fatalf("empty /model did not preserve current selection; model line count = %d, output = %q", count, got)
	}
}

func TestRunWithArgsConfigCommandShowsEffectiveNonSecretConfig(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "arena-rbx.json"), []byte(`{"arena":{"apiKeyEnv":"ARENA_TEST_SECRET","model":"arena/config-model"}}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldwd) })
	t.Setenv("ARENA_TEST_SECRET", "super-secret-value")

	in := strings.NewReader("/config\n/exit\n")
	var out bytes.Buffer
	if err := runWithArgs(context.Background(), in, &out, []string{"--model", "arena/flag-model"}); err != nil {
		t.Fatalf("runWithArgs() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, `"apiKeyEnv": "ARENA_TEST_SECRET"`) {
		t.Fatalf("config command missing api key env name: %q", got)
	}
	if !strings.Contains(got, `"model": "arena/flag-model"`) {
		t.Fatalf("config command missing effective model: %q", got)
	}
	if strings.Contains(got, "super-secret-value") {
		t.Fatalf("config command leaked secret: %q", got)
	}
}

func TestRunWithArgsModelsCommandListsStableUniqueDiscoveredModels(t *testing.T) {
	in := strings.NewReader("/models\n/exit\n")
	var out bytes.Buffer
	calls := 0

	listModels := func(context.Context) ([]string, error) {
		calls++
		return []string{" arena/model-b ", "arena/model-a", "arena/model-b", "", "   "}, nil
	}
	if err := runWithArgsAndModels(context.Background(), in, &out, nil, listModels); err != nil {
		t.Fatalf("runWithArgsAndModels() error = %v", err)
	}

	if calls != 1 {
		t.Fatalf("model discovery calls = %d, want 1", calls)
	}
	got := out.String()
	if !strings.Contains(got, "arena/model-a\narena/model-b\n") {
		t.Fatalf("models command missing stable unique discovered models: %q", got)
	}
	if strings.Contains(got, "arena/model-b \n") || strings.Count(got, "arena/model-b\n") != 1 {
		t.Fatalf("models command did not trim/deduplicate IDs: %q", got)
	}
}

func TestRunWithDependenciesDispatchesOrdinaryInputAsAgentTask(t *testing.T) {
	in := strings.NewReader("  inspect Workspace scripts  \n/exit\n")
	var out bytes.Buffer
	var tasks []string

	runTask := func(_ context.Context, task string) error {
		tasks = append(tasks, task)
		return nil
	}
	if err := runWithDependencies(context.Background(), in, &out, []string{"--model", "arena/test-model"}, nil, runTask); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	if len(tasks) != 1 || tasks[0] != "inspect Workspace scripts" {
		t.Fatalf("dispatched tasks = %#v, want one trimmed agent task", tasks)
	}
}

func TestRunWithDependenciesHistoryShowsSuccessfulSessionTasks(t *testing.T) {
	in := strings.NewReader("inspect Workspace scripts\n/history\n/exit\n")
	var out bytes.Buffer

	runTask := func(context.Context, string) error { return nil }
	if err := runWithDependencies(context.Background(), in, &out, []string{"--model", "arena/test-model"}, nil, runTask); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "task: inspect Workspace scripts\n") {
		t.Fatalf("history command missing successful session task: %q", got)
	}
}
