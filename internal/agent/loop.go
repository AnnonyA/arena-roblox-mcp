package agent

import (
	"context"
	"errors"
)

const DefaultMaxToolRounds = 12

var (
	ErrMaxToolRounds = errors.New("maximum tool rounds reached")
	ErrNoRoundFunc   = errors.New("tool loop round is not configured")
)

func RunToolLoop(ctx context.Context, maxRounds int, runRound func(context.Context) (bool, error)) error {
	if runRound == nil {
		return ErrNoRoundFunc
	}
	if ctx == nil {
		return ErrNoContext
	}
	if maxRounds <= 0 {
		maxRounds = DefaultMaxToolRounds
	}
	for round := 0; round < maxRounds; round++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		continueLoop, err := runRound(ctx)
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if !continueLoop {
			return nil
		}
	}
	return ErrMaxToolRounds
}
