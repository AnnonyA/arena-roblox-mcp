package agent

import (
	"context"
	"errors"
	"testing"
)

func TestRunToolLoopStopsAtConfiguredLimit(t *testing.T) {
	calls := 0
	err := RunToolLoop(context.Background(), 2, func(context.Context) (bool, error) {
		calls++
		return true, nil
	})

	if !errors.Is(err, ErrMaxToolRounds) {
		t.Fatalf("RunToolLoop error = %v, want %v", err, ErrMaxToolRounds)
	}
	if calls != 2 {
		t.Fatalf("round calls = %d, want 2", calls)
	}
}

func TestRunToolLoopReturnsWhenRoundFinishes(t *testing.T) {
	calls := 0
	err := RunToolLoop(context.Background(), 3, func(context.Context) (bool, error) {
		calls++
		return calls < 2, nil
	})

	if err != nil {
		t.Fatalf("RunToolLoop error = %v, want nil", err)
	}
	if calls != 2 {
		t.Fatalf("round calls = %d, want 2", calls)
	}
}

func TestRunToolLoopPropagatesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0

	err := RunToolLoop(ctx, DefaultMaxToolRounds, func(context.Context) (bool, error) {
		calls++
		return true, nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("RunToolLoop error = %v, want context canceled", err)
	}
	if calls != 0 {
		t.Fatalf("round calls = %d, want 0", calls)
	}
}
