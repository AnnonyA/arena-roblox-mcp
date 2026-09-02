package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

var (
	ErrNoConnector = errors.New("mcp connector is not configured")
	ErrNilSession  = errors.New("mcp connector returned a nil session")
)

type ToolResult struct {
	Content           json.RawMessage
	StructuredContent json.RawMessage
	IsError           bool
}

type Session interface {
	ListTools(context.Context) ([]Tool, error)
	CallTool(context.Context, string, json.RawMessage) (ToolResult, error)
	Close() error
}

type ConnectFunc func(context.Context) (Session, error)

type Client struct {
	mu       sync.Mutex
	connect  ConnectFunc
	session  Session
	inFlight chan struct{}
	registry *Registry
}

func NewClient(connect ConnectFunc) *Client {
	client := &Client{connect: connect}
	client.registry = NewRegistry(client.listTools)
	return client
}

func (c *Client) Tools(ctx context.Context) ([]Tool, error) {
	return c.registry.Tools(ctx)
}

func (c *Client) CallTool(ctx context.Context, name string, arguments json.RawMessage) (ToolResult, error) {
	session, err := c.ensureSession(ctx)
	if err != nil {
		return ToolResult{}, err
	}
	return session.CallTool(ctx, name, arguments)
}

func (c *Client) Close() error {
	for {
		c.mu.Lock()
		if c.inFlight != nil {
			inFlight := c.inFlight
			c.mu.Unlock()
			<-inFlight
			continue
		}

		session := c.session
		c.session = nil
		c.registry.Invalidate()
		c.mu.Unlock()

		if session == nil {
			return nil
		}
		return session.Close()
	}
}

func (c *Client) listTools(ctx context.Context) ([]Tool, error) {
	session, err := c.ensureSession(ctx)
	if err != nil {
		return nil, err
	}
	return session.ListTools(ctx)
}

func (c *Client) ensureSession(ctx context.Context) (Session, error) {
	for {
		c.mu.Lock()
		if c.session != nil {
			session := c.session
			c.mu.Unlock()
			return session, nil
		}

		if c.inFlight != nil {
			inFlight := c.inFlight
			c.mu.Unlock()

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-inFlight:
				continue
			}
		}

		if err := ctx.Err(); err != nil {
			c.mu.Unlock()
			return nil, err
		}

		if c.connect == nil {
			c.mu.Unlock()
			return nil, ErrNoConnector
		}

		inFlight := make(chan struct{})
		c.inFlight = inFlight
		c.mu.Unlock()

		session, err := c.connect(ctx)
		if err == nil && session == nil {
			err = ErrNilSession
		}

		c.mu.Lock()
		if err == nil {
			c.session = session
		}
		c.inFlight = nil
		close(inFlight)
		c.mu.Unlock()

		if err != nil {
			return nil, err
		}
		return session, nil
	}
}
