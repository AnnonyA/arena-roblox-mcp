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

	"github.com/AnnonyA/arena-roblox-mcp/internal/agent"
	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
	"github.com/AnnonyA/arena-roblox-mcp/internal/cli"
	"github.com/AnnonyA/arena-roblox-mcp/internal/config"
	"github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
	"github.com/AnnonyA/arena-roblox-mcp/internal/roblox"
	"github.com/AnnonyA/arena-roblox-mcp/internal/session"
)

const arenaBaseURL = "https://api.preview.arena.ai"
const sessionHistoryCapacity = 100
const conversationCapacity = 100

type listModelsFunc func(context.Context) ([]string, error)
type listToolsFunc func(context.Context) ([]string, error)
type listStudiosFunc func(context.Context) ([]roblox.StudioSession, error)
type runTaskFunc func(context.Context, []arena.Message) (string, error)

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
	return runWithDependencies(ctx, in, out, args, listModels, nil)
}

func runWithDependencies(ctx context.Context, in io.Reader, out io.Writer, args []string, listModels listModelsFunc, runTask runTaskFunc) error {
	return runWithToolDependencies(ctx, in, out, args, listModels, nil, runTask)
}

func runWithToolDependencies(ctx context.Context, in io.Reader, out io.Writer, args []string, listModels listModelsFunc, listTools listToolsFunc, runTask runTaskFunc) error {
	return runWithStudioDependencies(ctx, in, out, args, listModels, listTools, runTask, nil)
}

func runWithStudioDependencies(ctx context.Context, in io.Reader, out io.Writer, args []string, listModels listModelsFunc, listTools listToolsFunc, runTask runTaskFunc, listStudios listStudiosFunc) error {
	flags := flag.NewFlagSet("arena-rbx", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	model := flags.String("model", "", "Arena model ID")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			if out == nil {
				return errors.New("cli: nil output")
			}
			_, writeErr := io.WriteString(out, cli.HelpText())
			return writeErr
		}
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
		MCP:     "not connected",
		Studio:  "not connected",
		Model:   displayModel,
		Session: "default",
	}
	startup := strings.TrimSuffix(cli.StartupText(status), "> ")
	if _, err := io.WriteString(out, startup); err != nil {
		return err
	}

	history, err := session.NewHistory(sessionHistoryCapacity)
	if err != nil {
		return fmt.Errorf("create session history: %w", err)
	}
	conversation := agent.NewContext(conversationCapacity)

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

	var mcpClient *mcp.Client
	defer func() {
		if mcpClient != nil {
			_ = mcpClient.Close()
		}
	}()
	newMCPClient := func() (*mcp.Client, error) {
		if mcpClient != nil {
			return mcpClient, nil
		}
		server, ok := cfg.MCPServers["Roblox_Studio"]
		if !ok {
			return nil, errors.New("Roblox_Studio MCP server is not configured")
		}
		client, err := mcp.NewCommandClient(mcp.CommandConfig{Command: server.Command, Args: server.Args})
		if err != nil {
			return nil, err
		}
		mcpClient = client
		return mcpClient, nil
	}
	if listTools == nil {
		listTools = func(ctx context.Context) ([]string, error) {
			client, err := newMCPClient()
			if err != nil {
				return nil, err
			}
			tools, err := client.Tools(ctx)
			if err != nil {
				return nil, err
			}
			lines := make([]string, 0, len(tools))
			for _, tool := range tools {
				name := strings.TrimSpace(tool.Name)
				if name == "" {
					continue
				}
				description := strings.TrimSpace(tool.Description)
				if description == "" {
					lines = append(lines, name)
					continue
				}
				lines = append(lines, name+" — "+description)
			}
			sort.Strings(lines)
			return lines, nil
		}
	}
	if listStudios == nil {
		listStudios = func(ctx context.Context) ([]roblox.StudioSession, error) {
			client, err := newMCPClient()
			if err != nil {
				return nil, err
			}
			return roblox.DiscoverStudioSessions(ctx, client)
		}
	}

	if runTask == nil {
		var chatClient *arena.Client
		runTask = func(ctx context.Context, messages []arena.Message) (string, error) {
			modelID := strings.TrimSpace(cfg.Arena.Model)
			if modelID == "" {
				return "", errors.New("model not selected; use /model <id>")
			}
			if chatClient == nil {
				apiKey, err := config.ResolveAPIKey(cfg, os.Getenv)
				if err != nil {
					return "", err
				}
				chatClient = arena.NewClient(arena.ClientOptions{BaseURL: arenaBaseURL, APIKey: apiKey})
			}

			var writeErr error
			result, err := chatClient.StreamChat(ctx, arena.ChatRequest{
				Model:    modelID,
				Messages: messages,
			}, func(delta string) {
				if writeErr != nil {
					return
				}
				_, writeErr = io.WriteString(out, delta)
			})
			if err != nil {
				return "", err
			}
			if writeErr != nil {
				return "", writeErr
			}
			if result.Text != "" && !strings.HasSuffix(result.Text, "\n") {
				if _, err = io.WriteString(out, "\n"); err != nil {
					return "", err
				}
			}
			return result.Text, nil
		}
	}

	actions := cli.CommandActions{
		Clear: conversation.Clear,
		Status: func() cli.StartupStatus { return status },
		Models: listModels,
		Model: func(_ context.Context, id string) error {
			modelName := strings.TrimSpace(id)
			if modelName == "" {
				return errors.New("model ID is required")
			}
			cfg.Arena.Model = modelName
			status.Model = modelName
			return nil
		},
		Studio: func(ctx context.Context, id string) error {
			sessions, err := listStudios(ctx)
			if err != nil {
				return err
			}
			status.MCP = "connected"
			requestedID := strings.TrimSpace(id)
			if requestedID == "" && len(sessions) > 1 {
				sort.Slice(sessions, func(i, j int) bool { return sessions[i].ID < sessions[j].ID })
				if _, err := io.WriteString(out, "Multiple Roblox Studio sessions detected. Select one with /studio <studio_id>:\n"); err != nil {
					return err
				}
				lines := make([]string, 0, len(sessions))
				for _, studio := range sessions {
					line := studio.ID
					if studio.Name != "" {
						line += " — " + studio.Name
					}
					if studio.PlaceID != "" {
						line += " (place " + studio.PlaceID + ")"
					}
					lines = append(lines, line)
				}
				_, err = io.WriteString(out, strings.Join(lines, "\n")+"\n")
				return err
			}
			selected, err := roblox.SelectStudio(sessions, requestedID)
			if err != nil {
				return err
			}
			status.Studio = selected.ID
			return nil
		},
		Tools: func(ctx context.Context) ([]string, error) {
			tools, err := listTools(ctx)
			if err != nil {
				return nil, err
			}
			status.MCP = "connected"
			return tools, nil
		},
		History: func(context.Context) ([]string, error) {
			actions := history.Actions()
			lines := make([]string, 0, len(actions))
			for _, action := range actions {
				lines = append(lines, fmt.Sprintf("%s: %s", action.Tool, action.Summary))
			}
			return lines, nil
		},
		Config: func(context.Context) (string, error) {
			data, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return "", fmt.Errorf("render config: %w", err)
			}
			return string(data) + "\n", nil
		},
	}
	next := func(ctx context.Context, input cli.Input) (bool, error) {
		if input.Kind != cli.InputTask {
			return false, nil
		}
		if strings.TrimSpace(cfg.Arena.Model) == "" {
			models, err := listModels(ctx)
			if err != nil {
				return false, err
			}
			if _, err := io.WriteString(out, "Select a model with /model <id> before sending a task:\n"); err != nil {
				return false, err
			}
			if len(models) > 0 {
				if _, err := io.WriteString(out, strings.Join(models, "\n")+"\n"); err != nil {
					return false, err
				}
			}
			return false, nil
		}

		events := conversation.Events()
		messages := make([]arena.Message, 0, len(events)+1)
		for _, event := range events {
			messages = append(messages, arena.Message{Role: event.Role, Content: event.Content})
		}
		messages = append(messages, arena.Message{Role: "user", Content: input.Task})

		assistantReply, err := runTask(ctx, messages)
		if err != nil {
			return false, err
		}
		status.Arena = "connected"
		conversation.Add(agent.Event{Role: "user", Content: input.Task})
		if assistantReply != "" {
			conversation.Add(agent.Event{Role: "assistant", Content: assistantReply})
		}
		history.Add(session.Action{Tool: "task", Summary: input.Task})
		return false, nil
	}
	return cli.Run(ctx, in, out, cli.NewCommandHandlerWithActions(out, actions, next))
}
