package cli

import (
	"context"
	"errors"
	"io"
	"testing"
)

func TestCommandHandlerReportsShortWritesAcrossCommands(t *testing.T) {
	tests := []struct {
		name    string
		command string
		actions CommandActions
	}{
		{
			name:    "status",
			command: "status",
			actions: CommandActions{Status: func() StartupStatus { return StartupStatus{Arena: "connected"} }},
		},
		{
			name:    "models",
			command: "models",
			actions: CommandActions{Models: func(context.Context) ([]string, error) { return []string{"model-a"}, nil }},
		},
		{
			name:    "tools",
			command: "tools",
			actions: CommandActions{Tools: func(context.Context) ([]string, error) { return []string{"read_script"}, nil }},
		},
		{
			name:    "history",
			command: "history",
			actions: CommandActions{History: func(context.Context) ([]string, error) { return []string{"read_script: inspected"}, nil }},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewCommandHandlerWithActions(shortWriter{}, tt.actions, nil)
			_, err := handler(context.Background(), Input{Kind: InputCommand, Command: tt.command})
			if !errors.Is(err, io.ErrShortWrite) {
				t.Fatalf("%s error = %v, want %v", tt.command, err, io.ErrShortWrite)
			}
		})
	}
}
