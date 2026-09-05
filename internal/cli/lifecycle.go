package cli

import (
	"bufio"
	"context"
	"errors"
	"io"
)

type InputHandler func(context.Context, Input) (exit bool, err error)

func Run(ctx context.Context, in io.Reader, out io.Writer, handler InputHandler) error {
	if ctx == nil {
		return errors.New("cli: nil context")
	}
	if handler == nil {
		return errors.New("cli: nil input handler")
	}

	scanner := bufio.NewScanner(in)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, err := io.WriteString(out, "> "); err != nil {
			return err
		}
		if !scanner.Scan() {
			return scanner.Err()
		}

		input, err := ParseLine(scanner.Text())
		if err != nil {
			return err
		}
		if input.Kind == InputEmpty {
			continue
		}

		exit, err := handler(ctx, input)
		if err != nil {
			return err
		}
		if exit {
			return nil
		}
	}
}
