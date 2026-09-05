package cli

import (
	"context"
	"strings"
	"testing"
)

func TestRunRejectsNilOutput(t *testing.T) {
	err := Run(context.Background(), strings.NewReader("/exit\n"), nil, func(context.Context, Input) (bool, error) {
		return true, nil
	})
	if err == nil {
		t.Fatal("Run returned nil error for nil output")
	}
	if got, want := err.Error(), "cli: nil output"; got != want {
		t.Fatalf("Run error = %q, want %q", got, want)
	}
}
