package cli

import (
	"context"
	"errors"
	"io"
	"testing"
)

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return len(p) - 1, nil
}

func TestCommandHandlerReportsShortHelpWrite(t *testing.T) {
	handler := NewCommandHandler(shortWriter{}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "help"})
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("help error = %v, want %v", err, io.ErrShortWrite)
	}
}
