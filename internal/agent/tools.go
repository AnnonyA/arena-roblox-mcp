package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	mcppkg "github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
)

var (
	ErrUnknownTool          = errors.New("unknown tool")
	ErrInvalidToolArguments = errors.New("invalid tool arguments")
)

type ToolCaller interface {
	CallTool(context.Context, string, json.RawMessage) (mcppkg.ToolResult, error)
}

type ToolDispatcher struct {
	allowed map[string]struct{}
	caller  ToolCaller
}

func NewToolDispatcher(toolNames []string, caller ToolCaller) *ToolDispatcher {
	allowed := make(map[string]struct{}, len(toolNames))
	for _, name := range toolNames {
		allowed[name] = struct{}{}
	}
	return &ToolDispatcher{allowed: allowed, caller: caller}
}

func (d *ToolDispatcher) Dispatch(ctx context.Context, name string, arguments json.RawMessage) (mcppkg.ToolResult, error) {
	if _, ok := d.allowed[name]; !ok {
		return mcppkg.ToolResult{}, fmt.Errorf("%w: %s", ErrUnknownTool, name)
	}
	if !json.Valid(arguments) {
		return mcppkg.ToolResult{}, ErrInvalidToolArguments
	}
	return d.caller.CallTool(ctx, name, arguments)
}
