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

var ErrUnknownCommand = errors.New("unknown command")

type Input struct {
	Kind     InputKind
	Command  string
	Argument string
	Task     string
}

var supportedCommands = map[string]struct{}{
	"model":   {},
	"models":  {},
	"studio":  {},
	"status":  {},
	"tools":   {},
	"history": {},
	"diff":    {},
	"undo":    {},
	"clear":   {},
	"config":  {},
	"help":    {},
	"exit":    {},
}

var commandDescriptions = map[string]string{
	"model":   "change model",
	"models":  "list Arena models",
	"studio":  "select Studio instance",
	"status":  "show Arena/MCP/Studio state",
	"tools":   "show available MCP tools",
	"history": "show session actions and tool calls",
	"diff":    "show recorded changes",
	"undo":    "revert the latest supported reversible change",
	"clear":   "clear conversational context",
	"config":  "show effective non-secret configuration",
	"help":    "show commands",
	"exit":    "exit cleanly",
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
	command := strings.TrimPrefix(fields[0], "/")
	if _, ok := supportedCommands[command]; !ok {
		return Input{}, fmt.Errorf("%w: /%s", ErrUnknownCommand, command)
	}

	argument := strings.TrimSpace(strings.TrimPrefix(line, fields[0]))
	return Input{Kind: InputCommand, Command: command, Argument: argument}, nil
}

func HelpText() string {
	commands := make([]string, 0, len(supportedCommands))
	for command := range supportedCommands {
		commands = append(commands, command)
	}
	sort.Strings(commands)

	var help strings.Builder
	help.WriteString("Commands:\n")
	for _, command := range commands {
		fmt.Fprintf(&help, "/%s  %s\n", command, commandDescriptions[command])
	}
	return help.String()
}
