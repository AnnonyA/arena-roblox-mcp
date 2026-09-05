package cli

import (
	"errors"
	"testing"
)

func TestParseLineRoutesAgentTask(t *testing.T) {
	input, err := ParseLine("  inspect ServerScriptService.Main  ")
	if err != nil {
		t.Fatalf("ParseLine() error = %v", err)
	}
	if input.Kind != InputTask {
		t.Fatalf("kind = %v, want InputTask", input.Kind)
	}
	if input.Task != "inspect ServerScriptService.Main" {
		t.Fatalf("task = %q, want trimmed task", input.Task)
	}
}

func TestParseLineParsesSlashCommandAndArgument(t *testing.T) {
	input, err := ParseLine(" /model  arena-code ")
	if err != nil {
		t.Fatalf("ParseLine() error = %v", err)
	}
	if input.Kind != InputCommand {
		t.Fatalf("kind = %v, want InputCommand", input.Kind)
	}
	if input.Command != "model" {
		t.Fatalf("command = %q, want model", input.Command)
	}
	if input.Argument != "arena-code" {
		t.Fatalf("argument = %q, want arena-code", input.Argument)
	}
}

func TestParseLineRejectsUnknownSlashCommand(t *testing.T) {
	_, err := ParseLine("/explode")
	if !errors.Is(err, ErrUnknownCommand) {
		t.Fatalf("ParseLine() error = %v, want ErrUnknownCommand", err)
	}
}
