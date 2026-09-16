package cli

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestRunReportsShortPromptWrite(t *testing.T) {
	err := Run(context.Background(), strings.NewReader("/exit\n"), shortWriter{}, func(context.Context, Input) (bool, error) {
		return true, nil
	})
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("Run() error = %v, want %v", err, io.ErrShortWrite)
	}
}
