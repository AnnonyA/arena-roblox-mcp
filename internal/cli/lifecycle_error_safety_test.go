package cli

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRunSanitizesHandlerErrorsBeforeDisplay(t *testing.T) {
	in := strings.NewReader("fail task\n/exit\n")
	var out strings.Builder

	err := Run(context.Background(), in, &out, func(_ context.Context, input Input) (bool, error) {
		if input.Kind == InputTask {
			return false, errors.New("MCP failure: bad\x1b[31m\u200btool")
		}
		return input.Kind == InputCommand && input.Command == "exit", nil
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	const want = "> Error: MCP failure: bad[31mtool\n> "
	if got := out.String(); got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
