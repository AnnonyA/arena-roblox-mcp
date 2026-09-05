package cli

import (
	"context"
	"errors"
	"io"
)

type CommandActions struct {
	Clear  func()
	Status func() StartupStatus
}

func NewCommandHandler(out io.Writer, next InputHandler) InputHandler {
	return NewCommandHandlerWithActions(out, CommandActions{}, next)
}

func NewCommandHandlerWithClear(out io.Writer, clear func(), next InputHandler) InputHandler {
	return NewCommandHandlerWithActions(out, CommandActions{Clear: clear}, next)
}

func NewCommandHandlerWithActions(out io.Writer, actions CommandActions, next InputHandler) InputHandler {
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
		if input.Kind == InputCommand && input.Command == "clear" && actions.Clear != nil {
			actions.Clear()
			return false, nil
		}
		if input.Kind == InputCommand && input.Command == "status" && actions.Status != nil {
			if out == nil {
				return false, errors.New("cli: nil output")
			}
			_, err := io.WriteString(out, StatusText(actions.Status()))
			return false, err
		}
		if next == nil {
			return false, nil
		}
		return next(ctx, input)
	}
}
