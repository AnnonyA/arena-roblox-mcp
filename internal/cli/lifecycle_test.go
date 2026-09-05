package cli

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestRunDispatchesInputUntilHandlerExits(t *testing.T) {
	in := strings.NewReader("inspect Workspace\n/exit\n")
	var out strings.Builder
	var got []Input

	err := Run(context.Background(), in, &out, func(_ context.Context, input Input) (bool, error) {
		got = append(got, input)
		return input.Kind == InputCommand && input.Command == "exit", nil
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := []Input{
		{Kind: InputTask, Task: "inspect Workspace"},
		{Kind: InputCommand, Command: "exit"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("dispatched inputs = %#v, want %#v", got, want)
	}
	if gotPrompt := out.String(); gotPrompt != "> > " {
		t.Fatalf("prompt output = %q, want %q", gotPrompt, "> > ")
	}
}

func TestRunReportsLineErrorsAndContinues(t *testing.T) {
	in := strings.NewReader("/unknown\ninspect Workspace\nfail task\n/exit\n")
	var out strings.Builder
	var got []Input

	err := Run(context.Background(), in, &out, func(_ context.Context, input Input) (bool, error) {
		got = append(got, input)
		if input.Kind == InputTask && input.Task == "fail task" {
			return false, errors.New("agent unavailable")
		}
		return input.Kind == InputCommand && input.Command == "exit", nil
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := []Input{
		{Kind: InputTask, Task: "inspect Workspace"},
		{Kind: InputTask, Task: "fail task"},
		{Kind: InputCommand, Command: "exit"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("dispatched inputs = %#v, want %#v", got, want)
	}
	const wantOutput = "> Error: unknown command: /unknown\n> > Error: agent unavailable\n> "
	if gotOutput := out.String(); gotOutput != wantOutput {
		t.Fatalf("output = %q, want %q", gotOutput, wantOutput)
	}
}
