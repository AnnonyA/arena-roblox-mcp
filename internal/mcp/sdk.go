package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const maxToolListPages = 256

var (
	ErrMissingCommand      = errors.New("mcp command is required")
	ErrToolPaginationCycle = errors.New("mcp tool pagination cursor repeated")
	ErrToolPaginationLimit = errors.New("mcp tool pagination limit exceeded")
)

type CommandConfig struct {
	Command       string
	Args          []string
	ClientName    string
	ClientVersion string
}

func NewCommandClient(config CommandConfig) (*Client, error) {
	config = snapshotCommandConfig(config)
	if config.Command == "" {
		return nil, ErrMissingCommand
	}
	name := config.ClientName
	if name == "" {
		name = "arena-rbx"
	}
	version := config.ClientVersion
	if version == "" {
		version = "0.1.0"
	}

	connect := func(ctx context.Context) (Session, error) {
		client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: name, Version: version}, &sdkmcp.ClientOptions{
			Capabilities: &sdkmcp.ClientCapabilities{},
		})
		command := exec.Command(config.Command, config.Args...)
		session, err := client.Connect(ctx, &sdkmcp.CommandTransport{Command: command}, nil)
		if err != nil {
			return nil, err
		}
		return &sdkSession{session: session}, nil
	}

	return NewClient(connect), nil
}

func snapshotCommandConfig(config CommandConfig) CommandConfig {
	config.Command = strings.TrimSpace(config.Command)
	config.ClientName = strings.TrimSpace(config.ClientName)
	config.ClientVersion = strings.TrimSpace(config.ClientVersion)
	config.Args = append([]string(nil), config.Args...)
	return config
}

type sdkSession struct {
	session *sdkmcp.ClientSession
}

func (s *sdkSession) ListTools(ctx context.Context) ([]Tool, error) {
	var tools []Tool
	params := &sdkmcp.ListToolsParams{}
	seenCursors := make(map[string]struct{})
	page := 1
	for {
		result, err := s.session.ListTools(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, sdkTool := range result.Tools {
			if sdkTool == nil {
				continue
			}
			schema, err := json.Marshal(sdkTool.InputSchema)
			if err != nil {
				return nil, fmt.Errorf("marshal input schema for %q: %w", sdkTool.Name, err)
			}
			tools = append(tools, Tool{
				Name:        sdkTool.Name,
				Description: sdkTool.Description,
				InputSchema: schema,
			})
		}
		if result.NextCursor == "" {
			return tools, nil
		}
		page++
		if err := validateToolPageCursor(page, seenCursors, result.NextCursor); err != nil {
			return nil, err
		}
		params.Cursor = result.NextCursor
	}
}

func validateToolPageCursor(page int, seen map[string]struct{}, cursor string) error {
	if page > maxToolListPages {
		return fmt.Errorf("%w: maximum %d pages", ErrToolPaginationLimit, maxToolListPages)
	}
	if _, ok := seen[cursor]; ok {
		return fmt.Errorf("%w: %q", ErrToolPaginationCycle, cursor)
	}
	seen[cursor] = struct{}{}
	return nil
}

func (s *sdkSession) CallTool(ctx context.Context, name string, arguments json.RawMessage) (ToolResult, error) {
	args := map[string]any{}
	if len(arguments) != 0 {
		if err := json.Unmarshal(arguments, &args); err != nil {
			return ToolResult{}, fmt.Errorf("decode arguments for %q: %w", name, err)
		}
	}
	result, err := s.session.CallTool(ctx, &sdkmcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		return ToolResult{}, err
	}

	content, err := json.Marshal(result.Content)
	if err != nil {
		return ToolResult{}, fmt.Errorf("marshal result content for %q: %w", name, err)
	}
	var structured json.RawMessage
	if result.StructuredContent != nil {
		structured, err = json.Marshal(result.StructuredContent)
		if err != nil {
			return ToolResult{}, fmt.Errorf("marshal structured result for %q: %w", name, err)
		}
	}
	return ToolResult{Content: content, StructuredContent: structured, IsError: result.IsError}, nil
}

func (s *sdkSession) Close() error {
	return s.session.Close()
}
