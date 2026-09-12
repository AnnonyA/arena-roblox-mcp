package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrNoDiscoverer = errors.New("mcp tool discoverer is not configured")

type Tool struct {
	Name        string
	Description string
	InputSchema json.RawMessage
}
type DiscoverFunc func(context.Context) ([]Tool, error)
type registryAttempt struct {
	done chan struct{}
	err  error
}
type Registry struct {
	mu         sync.Mutex
	discover   DiscoverFunc
	tools      []Tool
	loaded     bool
	inFlight   *registryAttempt
	generation uint64
}

func NewRegistry(discover DiscoverFunc) *Registry { return &Registry{discover: discover} }
func (r *Registry) Tools(ctx context.Context) ([]Tool, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		r.mu.Lock()
		if err := ctx.Err(); err != nil {
			r.mu.Unlock()
			return nil, err
		}
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
			attempt := r.inFlight
			r.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-attempt.done:
				if attempt.err != nil {
					if ctx.Err() == nil && (errors.Is(attempt.err, context.Canceled) || errors.Is(attempt.err, context.DeadlineExceeded)) {
						continue
					}
					return nil, attempt.err
				}
				continue
			}
		}
		attempt := &registryAttempt{done: make(chan struct{})}
		r.inFlight = attempt
		generation := r.generation
		r.mu.Unlock()
		tools, err := r.discover(ctx)
		if err == nil {
			err = ctx.Err()
		}
		if err == nil {
			err = validateToolNames(tools)
		}
		r.mu.Lock()
		stale := generation != r.generation
		if err == nil && !stale {
			r.tools = cloneTools(tools)
			r.loaded = true
		}
		attempt.err = err
		r.inFlight = nil
		close(attempt.done)
		r.mu.Unlock()
		if err != nil {
			return nil, err
		}
		if stale {
			continue
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
func validateToolNames(tools []Tool) error {
	seen := make(map[string]struct{}, len(tools))
	for i, tool := range tools {
		trimmed := strings.TrimSpace(tool.Name)
		if trimmed == "" {
			return fmt.Errorf("blank MCP tool name at index %d", i)
		}
		if trimmed != tool.Name {
			return fmt.Errorf("MCP tool name %q at index %d has surrounding whitespace", tool.Name, i)
		}
		if _, ok := seen[tool.Name]; ok {
			return fmt.Errorf("duplicate MCP tool name %q at index %d", tool.Name, i)
		}
		seen[tool.Name] = struct{}{}
	}
	return nil
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
