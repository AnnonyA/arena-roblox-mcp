package cli

import (
	"context"
	"errors"
	"io"
	"testing"
)

type secondWriteShortWriter struct {
	writes int
}

func (w *secondWriteShortWriter) Write(p []byte) (int, error) {
	w.writes++
	if w.writes == 2 && len(p) > 0 {
		return len(p) - 1, nil
	}
	return len(p), nil
}

func TestModelSelectionPromptReportsFollowUpShortWrite(t *testing.T) {
	out := &secondWriteShortWriter{}
	handler := NewCommandHandlerWithActions(out, CommandActions{
		Models: func(context.Context) ([]string, error) {
			return []string{"model-a", "model-b"}, nil
		},
	}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "model"})
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("model selection error = %v, want %v", err, io.ErrShortWrite)
	}
	if out.writes != 2 {
		t.Fatalf("writes = %d, want 2", out.writes)
	}
}
