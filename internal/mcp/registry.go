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
