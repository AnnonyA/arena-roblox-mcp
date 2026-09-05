package cli

import (
	"bytes"
	"context"
	"testing"
)

func TestCommandHandlerWritesHelp(t *testing.T) {
	var out bytes.Buffer
	handler := NewCommandHandler(&out, nil)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "help"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("help unexpectedly requested exit")
	}
	if got, want := out.String(), HelpText(); got != want {
		t.Fatalf("help output = %q, want %q", got, want)
	}
}

func TestCommandHandlerRequestsExit(t *testing.T) {
	called := false
	next := func(context.Context, Input) (bool, error) {
		called = true
		return false, nil
	}
	handler := NewCommandHandler(nil, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "exit"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if !exit {
		t.Fatal("exit command did not request exit")
	}
	if called {
		t.Fatal("exit delegated to next handler")
	}
}

func TestCommandHandlerClearsConversation(t *testing.T) {
	cleared := false
	delegated := false
	next := func(context.Context, Input) (bool, error) {
		delegated = true
		return false, nil
	}
	handler := NewCommandHandlerWithClear(nil, func() { cleared = true }, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "clear"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("clear unexpectedly requested exit")
	}
	if !cleared {
		t.Fatal("clear command did not clear conversation")
	}
	if delegated {
		t.Fatal("clear delegated to next handler")
	}
}

func TestCommandHandlerWritesStatus(t *testing.T) {
	var out bytes.Buffer
	delegated := false
	next := func(context.Context, Input) (bool, error) {
		delegated = true
		return false, nil
	}
	status := func() StartupStatus {
		return StartupStatus{Arena: "connected", Studio: "offline", Model: "arena-code", Session: "default"}
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Status: status}, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "status"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("status unexpectedly requested exit")
	}
	if delegated {
		t.Fatal("status delegated to next handler")
	}
	const want = "Arena      connected\nStudio     offline\nModel      arena-code\nSession    default\n"
	if got := out.String(); got != want {
		t.Fatalf("status output = %q, want %q", got, want)
	}
}

func TestCommandHandlerWritesModels(t *testing.T) {
	var out bytes.Buffer
	delegated := false
	next := func(context.Context, Input) (bool, error) {
		delegated = true
		return false, nil
	}
	models := func(context.Context) ([]string, error) {
		return []string{"arena-code", "arena-fast"}, nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Models: models}, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "models"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("models unexpectedly requested exit")
	}
	if delegated {
		t.Fatal("models delegated to next handler")
	}
	const want = "arena-code\narena-fast\n"
	if got := out.String(); got != want {
		t.Fatalf("models output = %q, want %q", got, want)
	}
}

func TestCommandHandlerChangesModel(t *testing.T) {
	delegated := false
	next := func(context.Context, Input) (bool, error) {
		delegated = true
		return false, nil
	}
	var selected string
	model := func(_ context.Context, id string) error {
		selected = id
		return nil
	}
	handler := NewCommandHandlerWithActions(nil, CommandActions{Model: model}, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "model", Argument: "arena-code"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("model unexpectedly requested exit")
	}
	if delegated {
		t.Fatal("model delegated to next handler")
	}
	if selected != "arena-code" {
		t.Fatalf("selected model = %q, want %q", selected, "arena-code")
	}
}

func TestCommandHandlerWritesTools(t *testing.T) {
	var out bytes.Buffer
	delegated := false
	next := func(context.Context, Input) (bool, error) {
		delegated = true
		return false, nil
	}
	tools := func(context.Context) ([]string, error) {
		return []string{"read_script", "run_playtest"}, nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Tools: tools}, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "tools"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("tools unexpectedly requested exit")
	}
	if delegated {
		t.Fatal("tools delegated to next handler")
	}
	const want = "read_script\nrun_playtest\n"
	if got := out.String(); got != want {
		t.Fatalf("tools output = %q, want %q", got, want)
	}
}

func TestCommandHandlerSelectsStudio(t *testing.T) {
	delegated := false
	next := func(context.Context, Input) (bool, error) {
		delegated = true
		return false, nil
	}
	var selected string
	studio := func(_ context.Context, id string) error {
		selected = id
		return nil
	}
	handler := NewCommandHandlerWithActions(nil, CommandActions{Studio: studio}, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "studio", Argument: "studio-2"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("studio unexpectedly requested exit")
	}
	if delegated {
		t.Fatal("studio delegated to next handler")
	}
	if selected != "studio-2" {
		t.Fatalf("selected studio = %q, want %q", selected, "studio-2")
	}
}

func TestCommandHandlerWritesHistory(t *testing.T) {
	var out bytes.Buffer
	delegated := false
	next := func(context.Context, Input) (bool, error) {
		delegated = true
		return false, nil
	}
	history := func(context.Context) ([]string, error) {
		return []string{"read_script: inspected PlayerController", "write_script: updated jump logic"}, nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{History: history}, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "history"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("history unexpectedly requested exit")
	}
	if delegated {
		t.Fatal("history delegated to next handler")
	}
	const want = "read_script: inspected PlayerController\nwrite_script: updated jump logic\n"
	if got := out.String(); got != want {
		t.Fatalf("history output = %q, want %q", got, want)
	}
}

func TestCommandHandlerWritesDiff(t *testing.T) {
	var out bytes.Buffer
	delegated := false
	next := func(context.Context, Input) (bool, error) {
		delegated = true
		return false, nil
	}
	diff := func(context.Context) (string, error) {
		return "--- before\n+++ after\n-old\n+new\n", nil
	}
	handler := NewCommandHandlerWithActions(&out, CommandActions{Diff: diff}, next)

	exit, err := handler(context.Background(), Input{Kind: InputCommand, Command: "diff"})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if exit {
		t.Fatal("diff unexpectedly requested exit")
	}
	if delegated {
		t.Fatal("diff delegated to next handler")
	}
	const want = "--- before\n+++ after\n-old\n+new\n"
	if got := out.String(); got != want {
		t.Fatalf("diff output = %q, want %q", got, want)
	}
}
