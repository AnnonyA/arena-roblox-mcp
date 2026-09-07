package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

func TestFakeEndToEndModelSelectionTasksHistoryAndStatus(t *testing.T) {
	in := strings.NewReader(strings.Join([]string{
		"/models",
		"/model arena/model-a",
		"inspect Workspace scripts",
		"fix the issue",
		"/history",
		"/status",
		"/exit",
		"",
	}, "\n"))
	var out bytes.Buffer

	modelCalls := 0
	listModels := func(context.Context) ([]string, error) {
		modelCalls++
		return []string{"arena/model-b", " arena/model-a ", "arena/model-b"}, nil
	}

	var requests [][]arena.Message
	runTask := func(_ context.Context, messages []arena.Message) (string, error) {
		copied := append([]arena.Message(nil), messages...)
		requests = append(requests, copied)
		if len(requests) == 1 {
			return "inspection complete", nil
		}
		return "fixed", nil
	}

	if err := runWithDependencies(context.Background(), in, &out, nil, listModels, runTask); err != nil {
		t.Fatalf("runWithDependencies() error = %v", err)
	}

	if modelCalls != 2 {
		t.Fatalf("model discovery calls = %d, want 2 (/models and /model validation)", modelCalls)
	}
	if len(requests) != 2 {
		t.Fatalf("task requests = %d, want 2", len(requests))
	}

	first := requests[0]
	if len(first) != 1 || first[0].Role != "user" || first[0].Content != "inspect Workspace scripts" {
		t.Fatalf("first request = %#v, want current user task", first)
	}
	second := requests[1]
	if len(second) != 3 {
		t.Fatalf("second request = %#v, want prior user/assistant context plus current task", second)
	}
	wantSecond := []arena.Message{
		{Role: "user", Content: "inspect Workspace scripts"},
		{Role: "assistant", Content: "inspection complete"},
		{Role: "user", Content: "fix the issue"},
	}
	for i := range wantSecond {
		if second[i] != wantSecond[i] {
			t.Fatalf("second request[%d] = %#v, want %#v", i, second[i], wantSecond[i])
		}
	}

	got := out.String()
	for _, want := range []string{
		"arena/model-a\narena/model-b\n",
		"task: inspect Workspace scripts\n",
		"task: fix the issue\n",
		"Arena      connected\n",
		"Model      arena/model-a\n",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q: %q", want, got)
		}
	}
}
