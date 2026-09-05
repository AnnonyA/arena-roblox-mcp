package cli

import (
	"context"
	"errors"
	"io"
	"strings"
)

type CommandActions struct {
	Clear   func()
	Status  func() StartupStatus
	Models  func(context.Context) ([]string, error)
	Model   func(context.Context, string) error
	Studio  func(context.Context, string) error
	Tools   func(context.Context) ([]string, error)
	History func(context.Context) ([]string, error)
	Diff    func(context.Context) (string, error)
	Undo    func(context.Context) error
	Config  func(context.Context) (string, error)
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
		if input.Kind == InputCommand && input.Command == "models" && actions.Models != nil {
			models, err := actions.Models(ctx)
			if err != nil {
				return false, err
			}
			if out == nil {
				return false, errors.New("cli: nil output")
			}
			if len(models) == 0 {
				return false, nil
			}
			_, err = io.WriteString(out, strings.Join(models, "\n")+"\n")
			return false, err
		}
		if input.Kind == InputCommand && input.Command == "model" && actions.Model != nil {
			return false, actions.Model(ctx, input.Argument)
		}
		if input.Kind == InputCommand && input.Command == "studio" && actions.Studio != nil {
			return false, actions.Studio(ctx, input.Argument)
		}
		if input.Kind == InputCommand && input.Command == "tools" && actions.Tools != nil {
			tools, err := actions.Tools(ctx)
			if err != nil {
				return false, err
			}
			if out == nil {
				return false, errors.New("cli: nil output")
			}
			if len(tools) == 0 {
				return false, nil
			}
			_, err = io.WriteString(out, strings.Join(tools, "\n")+"\n")
			return false, err
		}
		if input.Kind == InputCommand && input.Command == "history" && actions.History != nil {
			history, err := actions.History(ctx)
			if err != nil {
				return false, err
			}
			if out == nil {
				return false, errors.New("cli: nil output")
			}
			if len(history) == 0 {
				return false, nil
			}
			_, err = io.WriteString(out, strings.Join(history, "\n")+"\n")
			return false, err
		}
		if input.Kind == InputCommand && input.Command == "diff" && actions.Diff != nil {
			diff, err := actions.Diff(ctx)
			if err != nil {
				return false, err
			}
			if out == nil {
				return false, errors.New("cli: nil output")
			}
			if diff == "" {
				return false, nil
			}
			_, err = io.WriteString(out, diff)
			return false, err
		}
		if input.Kind == InputCommand && input.Command == "undo" && actions.Undo != nil {
			return false, actions.Undo(ctx)
		}
		if input.Kind == InputCommand && input.Command == "config" && actions.Config != nil {
			config, err := actions.Config(ctx)
			if err != nil {
				return false, err
			}
			if out == nil {
				return false, errors.New("cli: nil output")
			}
			if config == "" {
				return false, nil
			}
			_, err = io.WriteString(out, config)
			return false, err
		}
		if next == nil {
			return false, nil
		}
		return next(ctx, input)
	}
}
