package cli

import (
	"context"
	"errors"
	"io"
	"testing"
)

func TestCommandHandlerReportsShortHelpWrite(t *testing.T) {
	handler := NewCommandHandler(shortWriter{}, nil)

	_, err := handler(context.Background(), Input{Kind: InputCommand, Command: "help"})
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("help error = %v, want %v", err, io.ErrShortWrite)
	}
}
