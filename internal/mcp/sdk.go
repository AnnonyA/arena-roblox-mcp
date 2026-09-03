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

var ErrMissingCommand = errors.New("mcp command is required")

type CommandConfig struct {
	Command       string
	Args          []string
	ClientName    string
	ClientVersion string
}

func NewCommandClient(config CommandConfig) (*Client, error) {
	if strings.TrimSpace(config.Command) == "" {
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

type sdkSession struct {
	session *sdkmcp.ClientSession
}

func (s *sdkSession) ListTools(ctx context.Context) ([]Tool, error) {
	var tools []Tool
	params := &sdkmcp.ListToolsParams{}
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
		params.Cursor = result.NextCursor
	}
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
