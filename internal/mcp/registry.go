package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

var ErrNoDiscoverer = errors.New("mcp tool discoverer is not configured")

type Tool struct {
	Name        string
	Description string
	InputSchema json.RawMessage
}

type DiscoverFunc func(context.Context) ([]Tool, error)

type Registry struct {
	mu         sync.Mutex
	discover   DiscoverFunc
	tools      []Tool
	loaded     bool
	inFlight   chan struct{}
	generation uint64
}

func NewRegistry(discover DiscoverFunc) *Registry { return &Registry{discover: discover} }

func (r *Registry) Tools(ctx context.Context) ([]Tool, error) {
	for {
		r.mu.Lock()
		if r.loaded {
			tools := cloneTools(r.tools)
			r.mu.Unlock()
			return tools, nil
		}
		if r.discover == nil {
			r.mu.Unlock()
			return nil, ErrNoDiscoverer
		}

		if r.inFlight != nil {
			inFlight := r.inFlight
			r.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-inFlight:
				continue
			}
		}
		inFlight := make(chan struct{})
		r.inFlight = inFlight
		generation := r.generation
		r.mu.Unlock()
		tools, err := r.discover(ctx)
		r.mu.Lock()
		if err == nil && generation == r.generation {
			r.tools = cloneTools(tools)
			r.loaded = true
		}
		r.inFlight = nil
		close(inFlight)
		r.mu.Unlock()
		if err != nil {
			return nil, err
		}
		return cloneTools(tools), nil
	}
}

func (r *Registry) Invalidate() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.generation++
	r.tools = nil
	r.loaded = false
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
