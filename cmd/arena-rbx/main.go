package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/AnnonyA/arena-roblox-mcp/internal/cli"
	"github.com/AnnonyA/arena-roblox-mcp/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := runWithArgs(ctx, os.Stdin, os.Stdout, os.Args[1:]); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, "arena-rbx:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, in io.Reader, out io.Writer) error {
	return runWithArgs(ctx, in, out, nil)
}

func runWithArgs(ctx context.Context, in io.Reader, out io.Writer, args []string) error {
	flags := flag.NewFlagSet("arena-rbx", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	model := flags.String("model", "", "Arena model ID")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	cfg, err := config.Load("arena-rbx.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	modelName := strings.TrimSpace(*model)
	if modelName == "" {
		modelName = strings.TrimSpace(cfg.Arena.Model)
	}
	cfg.Arena.Model = modelName

	displayModel := modelName
	if displayModel == "" {
		displayModel = "not selected"
	}
	status := cli.StartupStatus{
		Arena:   "not connected",
		Studio:  "not connected",
		Model:   displayModel,
		Session: "default",
	}
	startup := strings.TrimSuffix(cli.StartupText(status), "> ")
	if _, err := io.WriteString(out, startup); err != nil {
		return err
	}

	actions := cli.CommandActions{
		Status: func() cli.StartupStatus { return status },
		Config: func(context.Context) (string, error) {
			data, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return "", fmt.Errorf("render config: %w", err)
			}
			return string(data) + "\n", nil
		},
	}
	return cli.Run(ctx, in, out, cli.NewCommandHandlerWithActions(out, actions, nil))
}
