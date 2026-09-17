package cli

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
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
		displayCommand := safeCommandDisplay(command)
		if suggestion := suggestedCommand(command); suggestion != "" {
			return Input{}, fmt.Errorf("%w: /%s (did you mean /%s? try /help)", ErrUnknownCommand, displayCommand, suggestion)
		}
		return Input{}, fmt.Errorf("%w: /%s (try /help)", ErrUnknownCommand, displayCommand)
	}

	argument := strings.TrimSpace(strings.TrimPrefix(line, fields[0]))
	if argument != "" && metadata.argumentHint == "" {
		return Input{}, fmt.Errorf("%w: /%s (usage: %s)", ErrUnexpectedCommandArgument, command, commandUsage(command, metadata))
	}
	return Input{Kind: InputCommand, Command: command, Argument: argument}, nil
}

func safeCommandDisplay(command string) string {
	quoted := strconv.Quote(command)
	return quoted[1 : len(quoted)-1]
}

func suggestedCommand(command string) string {
	const maxDistance = 2

	best := ""
	bestDistance := maxDistance + 1
	for candidate := range commands {
		distance := editDistance(command, candidate)
		if distance < bestDistance || distance == bestDistance && candidate < best {
			best = candidate
			bestDistance = distance
		}
	}
	if bestDistance > maxDistance {
		return ""
	}
	return best
}

func editDistance(a, b string) int {
	left, right := []rune(a), []rune(b)
	previous := make([]int, len(right)+1)
	current := make([]int, len(right)+1)
	for j := range previous {
		previous[j] = j
	}
	for i, leftRune := range left {
		current[0] = i + 1
		for j, rightRune := range right {
			cost := 0
			if leftRune != rightRune {
				cost = 1
			}
			current[j+1] = min3(current[j]+1, previous[j+1]+1, previous[j]+cost)
		}
		previous, current = current, previous
	}
	return previous[len(right)]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func commandUsage(command string, metadata commandMetadata) string {
	usage := "/" + command
	if metadata.argumentHint != "" {
		usage += " [" + metadata.argumentHint + "]"
	}
	return usage
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
		fmt.Fprintf(&help, "%-14s  %s\n", commandUsage(command, metadata), metadata.description)
	}
	return help.String()
}
