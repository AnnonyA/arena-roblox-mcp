package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type InputKind uint8

const (
	InputEmpty InputKind = iota
	InputCommand
	InputTask
)

var (
	ErrUnknownCommand            = errors.New("unknown command")
	ErrUnexpectedCommandArgument = errors.New("unexpected command argument")
)

type Input struct {
	Kind     InputKind
	Command  string
	Argument string
	Task     string
}

type commandMetadata struct {
	description  string
	argumentHint string
}

var commands = map[string]commandMetadata{
	"model":   {description: "change model", argumentHint: "id"},
	"models":  {description: "list Arena models"},
	"studio":  {description: "select Studio instance", argumentHint: "id"},
	"status":  {description: "show Arena/MCP/Studio state"},
	"tools":   {description: "show available MCP tools"},
	"history": {description: "show session actions and tool calls"},
	"diff":    {description: "show recorded changes"},
	"undo":    {description: "revert the latest supported reversible change"},
	"clear":   {description: "clear conversational context"},
	"config":  {description: "show effective non-secret configuration"},
	"help":    {description: "show commands"},
	"exit":    {description: "exit cleanly"},
}

func ParseLine(line string) (Input, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return Input{Kind: InputEmpty}, nil
	}
	if !strings.HasPrefix(line, "/") {
		return Input{Kind: InputTask, Task: line}, nil
	}

	fields := strings.Fields(line)
	command := strings.ToLower(strings.TrimPrefix(fields[0], "/"))
	metadata, ok := commands[command]
	if !ok {
		return Input{}, fmt.Errorf("%w: /%s", ErrUnknownCommand, command)
	}

	argument := strings.TrimSpace(strings.TrimPrefix(line, fields[0]))
	if argument != "" && metadata.argumentHint == "" {
		return Input{}, fmt.Errorf("%w: /%s", ErrUnexpectedCommandArgument, command)
	}
	return Input{Kind: InputCommand, Command: command, Argument: argument}, nil
}

func HelpText() string {
	commandNames := make([]string, 0, len(commands))
	for command := range commands {
		commandNames = append(commandNames, command)
	}
	sort.Strings(commandNames)

	var help strings.Builder
	help.WriteString("Commands:\n")
	for _, command := range commandNames {
		metadata := commands[command]
		usage := "/" + command
		if metadata.argumentHint != "" {
			usage += " [" + metadata.argumentHint + "]"
		}
		fmt.Fprintf(&help, "%-14s  %s\n", usage, metadata.description)
	}
	return help.String()
}
