package cli

import (
	"context"
	"errors"
	"io"
)

func NewCommandHandler(out io.Writer, next InputHandler) InputHandler {
	return NewCommandHandlerWithClear(out, nil, next)
}

func NewCommandHandlerWithClear(out io.Writer, clear func(), next InputHandler) InputHandler {
	return func(ctx context.Context, input Input) (bool, error) {
		if input.Kind == InputCommand && input.Command == "help" {
			if out == nil {
				return false, errors.New("cli: nil output")
			}
			_, err := io.WriteString(out, HelpText())
			return false, err
		}
		if input.Kind == InputCommand && input.Command == "exit" {
			return true, nil
		}
		if input.Kind == InputCommand && input.Command == "clear" && clear != nil {
			clear()
			return false, nil
		}
		if next == nil {
			return false, nil
		}
		return next(ctx, input)
	}
}
