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
	"sort"
	"strings"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
	"github.com/AnnonyA/arena-roblox-mcp/internal/cli"
	"github.com/AnnonyA/arena-roblox-mcp/internal/config"
)

const arenaBaseURL = "https://api.preview.arena.ai"

type listModelsFunc func(context.Context) ([]string, error)

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
	return runWithArgsAndModels(ctx, in, out, args, nil)
}

func runWithArgsAndModels(ctx context.Context, in io.Reader, out io.Writer, args []string, listModels listModelsFunc) error {
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

	if listModels == nil {
		listModels = func(ctx context.Context) ([]string, error) {
			apiKey, err := config.ResolveAPIKey(cfg, os.Getenv)
			if err != nil {
				return nil, err
			}
			client := arena.NewClient(arena.ClientOptions{BaseURL: arenaBaseURL, APIKey: apiKey})
			models, err := client.ListModels(ctx)
			if err != nil {
				return nil, err
			}
			ids := make([]string, 0, len(models))
			for _, model := range models {
				if id := strings.TrimSpace(model.ID); id != "" {
					ids = append(ids, id)
				}
			}
			return ids, nil
		}
	}

	baseListModels := listModels
	listModels = func(ctx context.Context) ([]string, error) {
		models, err := baseListModels(ctx)
		if err != nil {
			return nil, err
		}
		unique := make(map[string]struct{}, len(models))
		for _, model := range models {
			if id := strings.TrimSpace(model); id != "" {
				unique[id] = struct{}{}
			}
		}
		models = models[:0]
		for id := range unique {
			models = append(models, id)
		}
		sort.Strings(models)
		return models, nil
	}

	actions := cli.CommandActions{
		Status: func() cli.StartupStatus { return status },
		Models: listModels,
		Model: func(_ context.Context, id string) error {
			modelName := strings.TrimSpace(id)
			cfg.Arena.Model = modelName
			status.Model = modelName
			if status.Model == "" {
				status.Model = "not selected"
			}
			return nil
		},
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
