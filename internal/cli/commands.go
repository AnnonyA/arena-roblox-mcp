package cli

import (
	"errors"
	"fmt"
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
