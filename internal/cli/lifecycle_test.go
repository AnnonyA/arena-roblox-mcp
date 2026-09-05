package cli

import (
	"context"
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
