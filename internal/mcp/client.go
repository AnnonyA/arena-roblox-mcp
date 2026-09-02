package mcp

import (
	"context"
	"encoding/json"
	"sync"
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
	c.mu.Lock()
	session := c.session
	c.session = nil
	c.registry.Invalidate()
	c.mu.Unlock()

	if session == nil {
		return nil
	}
	return session.Close()
}

func (c *Client) listTools(ctx context.Context) ([]Tool, error) {
	session, err := c.ensureSession(ctx)
	if err != nil {
		return nil, err
	}
	return session.ListTools(ctx)
}

func (c *Client) ensureSession(ctx context.Context) (Session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.session != nil {
		return c.session, nil
	}

	session, err := c.connect(ctx)
	if err != nil {
		return nil, err
	}
	c.session = session
	return session, nil
}
