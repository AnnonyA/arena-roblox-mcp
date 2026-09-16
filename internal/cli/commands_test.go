package cli

import (
	"errors"
	"strings"
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

func TestParseLineAcceptsCaseInsensitiveSlashCommand(t *testing.T) {
	input, err := ParseLine(" /MoDeL  arena-code ")
	if err != nil {
		t.Fatalf("ParseLine() error = %v", err)
	}
	if input.Command != "model" {
		t.Fatalf("command = %q, want normalized model", input.Command)
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

func TestParseLineRejectsArgumentsForArgumentlessCommand(t *testing.T) {
	_, err := ParseLine("/status unexpected")
	if !errors.Is(err, ErrUnexpectedCommandArgument) {
		t.Fatalf("ParseLine() error = %v, want ErrUnexpectedCommandArgument", err)
	}
}

func TestCommandMetadataMatchesSupportedCommands(t *testing.T) {
	for command := range supportedCommands {
		description, ok := commandDescriptions[command]
		if !ok || strings.TrimSpace(description) == "" {
			t.Errorf("supported command %q is missing a description", command)
		}
	}
	for command := range commandDescriptions {
		if _, ok := supportedCommands[command]; !ok {
			t.Errorf("description exists for unsupported command %q", command)
		}
	}
	for command := range commandsWithArguments {
		if _, ok := supportedCommands[command]; !ok {
			t.Errorf("argument support exists for unsupported command %q", command)
		}
		hint, ok := commandArgumentHints[command]
		if !ok || strings.TrimSpace(hint) == "" {
			t.Errorf("command %q accepts arguments but has no help hint", command)
		}
	}
	for command := range commandArgumentHints {
		if _, ok := commandsWithArguments[command]; !ok {
			t.Errorf("argument hint exists for command %q that does not accept arguments", command)
		}
	}
}

func TestHelpTextListsEverySupportedCommandWithDescription(t *testing.T) {
	help := HelpText()
	for command := range supportedCommands {
		needle := "/" + command
		if !strings.Contains(help, needle) {
			t.Fatalf("HelpText() missing command description for %q", command)
		}
	}
}

func TestHelpTextShowsOptionalArguments(t *testing.T) {
	help := HelpText()
	for _, want := range []string{"/model [id]", "/studio [id]"} {
		if !strings.Contains(help, want) {
			t.Fatalf("HelpText() = %q, want %q", help, want)
		}
	}
	if strings.Contains(help, "/status [") {
		t.Fatalf("HelpText() marks argumentless /status as accepting an argument: %q", help)
	}
}
