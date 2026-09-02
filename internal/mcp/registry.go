package mcp

import (
	"context"
	"encoding/json"
	"sync"
)

type Tool struct {
	Name        string
	Description string
	InputSchema json.RawMessage
}

type DiscoverFunc func(context.Context) ([]Tool, error)

type Registry struct {
	mu       sync.Mutex
	discover DiscoverFunc
	tools    []Tool
	loaded   bool
}

func NewRegistry(discover DiscoverFunc) *Registry {
	return &Registry{discover: discover}
}

func (r *Registry) Tools(ctx context.Context) ([]Tool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.loaded {
		return cloneTools(r.tools), nil
	}

	tools, err := r.discover(ctx)
	if err != nil {
		return nil, err
	}

	r.tools = cloneTools(tools)
	r.loaded = true
	return cloneTools(r.tools), nil
}

func cloneTools(tools []Tool) []Tool {
	if tools == nil {
		return nil
	}
	cloned := make([]Tool, len(tools))
	for i, tool := range tools {
		cloned[i] = tool
		cloned[i].InputSchema = append(json.RawMessage(nil), tool.InputSchema...)
	}
	return cloned
}
