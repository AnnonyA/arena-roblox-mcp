package agent

import (
	"context"
	"errors"
	"testing"
)

func TestRunToolLoopRejectsNilContext(t *testing.T) {
	err := RunToolLoop(nil, 1, func(ctx context.Context) (bool, error) {
		t.Fatal("runRound should not be called with a nil context")
		return false, nil
	})
	if !errors.Is(err, ErrNoContext) {
		t.Fatalf("RunToolLoop(nil) error = %v, want ErrNoContext", err)
	}
}
