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

func TestParseLineUnknownSlashCommandSuggestsHelp(t *testing.T) {
	_, err := ParseLine("/explode")
	if err == nil {
		t.Fatal("ParseLine() error = nil, want unknown command error")
	}
	if !strings.Contains(err.Error(), "try /help") {
		t.Fatalf("ParseLine() error = %q, want actionable /help hint", err)
	}
}

func TestParseLineUnknownSlashCommandSuggestsCloseMatch(t *testing.T) {
	_, err := ParseLine("/stats")
	if err == nil {
		t.Fatal("ParseLine() error = nil, want unknown command error")
	}
	if !strings.Contains(err.Error(), "did you mean /status?") {
		t.Fatalf("ParseLine() error = %q, want /status typo hint", err)
	}
}

func TestParseLineUnknownSlashCommandAvoidsDistantSuggestion(t *testing.T) {
	_, err := ParseLine("/explode")
	if err == nil {
		t.Fatal("ParseLine() error = nil, want unknown command error")
	}
	if strings.Contains(err.Error(), "did you mean") {
		t.Fatalf("ParseLine() error = %q, want no low-confidence typo hint", err)
	}
}

func TestParseLineUnknownSlashCommandEscapesControlCharacters(t *testing.T) {
	_, err := ParseLine("/bad\x1b[31m")
	if err == nil {
		t.Fatal("ParseLine() error = nil, want unknown command error")
	}
	if strings.ContainsRune(err.Error(), '\x1b') {
		t.Fatalf("ParseLine() error contains terminal escape: %q", err)
	}
	if !strings.Contains(err.Error(), `/bad\x1b[31m`) {
		t.Fatalf("ParseLine() error = %q, want escaped control character", err)
	}
}

func TestParseLineUnknownSlashCommandBoundsDisplayedInput(t *testing.T) {
	_, err := ParseLine("/" + strings.Repeat("x", 4096))
	if err == nil {
		t.Fatal("ParseLine() error = nil, want unknown command error")
	}
	if len(err.Error()) > 256 {
		t.Fatalf("ParseLine() error length = %d, want bounded diagnostic", len(err.Error()))
	}
	if !strings.Contains(err.Error(), "...") {
		t.Fatalf("ParseLine() error = %q, want truncation marker", err)
	}
}

func TestEditDistanceAtMostRejectsLengthDifferenceBeyondLimit(t *testing.T) {
	if _, ok := editDistanceAtMost(strings.Repeat("x", 4096), "status", 2); ok {
		t.Fatal("editDistanceAtMost() ok = true, want immediate rejection for distant lengths")
	}
}

func TestEditDistanceAtMostFindsCloseUnicodeMatch(t *testing.T) {
	distance, ok := editDistanceAtMost("státus", "status", 2)
	if !ok || distance != 1 {
		t.Fatalf("editDistanceAtMost() = (%d, %v), want (1, true)", distance, ok)
	}
}

func TestParseLineRejectsArgumentsForArgumentlessCommand(t *testing.T) {
	_, err := ParseLine("/status unexpected")
	if !errors.Is(err, ErrUnexpectedCommandArgument) {
		t.Fatalf("ParseLine() error = %v, want ErrUnexpectedCommandArgument", err)
	}
}

func TestParseLineUnexpectedArgumentShowsUsage(t *testing.T) {
	_, err := ParseLine("/status unexpected")
	if err == nil {
		t.Fatal("ParseLine() error = nil, want unexpected argument error")
	}
	if !strings.Contains(err.Error(), "usage: /status") {
		t.Fatalf("ParseLine() error = %q, want actionable usage hint", err)
	}
}

func TestCommandMetadataIsComplete(t *testing.T) {
	for command, metadata := range commands {
		if strings.TrimSpace(command) == "" {
			t.Error("command metadata contains an empty command name")
		}
		if strings.TrimSpace(metadata.description) == "" {
			t.Errorf("command %q is missing a description", command)
		}
	}
}

func TestHelpTextRendersEveryCommandMetadataEntry(t *testing.T) {
	help := HelpText()
	for command, metadata := range commands {
		usage := "/" + command
		if metadata.argumentHint != "" {
			usage += " [" + metadata.argumentHint + "]"
		}
		if !strings.Contains(help, usage) {
			t.Errorf("HelpText() missing usage %q", usage)
		}
		if !strings.Contains(help, metadata.description) {
			t.Errorf("HelpText() missing description %q for /%s", metadata.description, command)
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
