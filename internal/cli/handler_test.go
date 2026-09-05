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
