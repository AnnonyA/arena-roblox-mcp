package agent

import (
	"context"
	"errors"
	"testing"
)

func TestRunToolLoopRejectsNilRound(t *testing.T) {
	err := RunToolLoop(context.Background(), 1, nil)
	if !errors.Is(err, ErrNoRoundFunc) {
		t.Fatalf("RunToolLoop(nil) error = %v, want ErrNoRoundFunc", err)
	}
}
