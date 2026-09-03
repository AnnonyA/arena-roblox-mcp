package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

var (
	ErrNoConnector  = errors.New("mcp connector is not configured")
	ErrNilSession   = errors.New("mcp connector returned a nil session")
	ErrClientClosed = errors.New("mcp client is closed")
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

type connectAttempt struct {
	done   chan struct{}
	cancel context.CancelFunc
	err    error
}

type Client struct {
	mu       sync.Mutex
	connect  ConnectFunc
	session  Session
	inFlight *connectAttempt
	closed   bool
	registry *Registry
}

func NewClient(connect ConnectFunc) *Client {
	client := &Client{connect: connect}
	client.registry = NewRegistry(client.listTools)
	return client
}

func (c *Client) Tools(ctx context.Context) ([]Tool, error) { return c.registry.Tools(ctx) }

func (c *Client) CallTool(ctx context.Context, name string, arguments json.RawMessage) (ToolResult, error) {
	if err := ctx.Err(); err != nil {
		return ToolResult{}, err
	}
	session, err := c.ensureSession(ctx)
	if err != nil {
		return ToolResult{}, err
	}
	result, err := session.CallTool(ctx, name, arguments)
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ToolResult{}, ctxErr
	}
	if err != nil {
		return ToolResult{}, err
	}
	return result, nil
}

func (c *Client) Close() error { return c.CloseContext(context.Background()) }

func (c *Client) CloseContext(ctx context.Context) error {
	for {
		c.mu.Lock()
		if err := ctx.Err(); err != nil {
			c.mu.Unlock()
			return err
		}
		c.closed = true
		if c.inFlight != nil {
			attempt := c.inFlight
			attempt.cancel()
			c.mu.Unlock()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-attempt.done:
				continue
			}
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
		if c.closed {
			c.mu.Unlock()
			return nil, ErrClientClosed
		}
		if err := ctx.Err(); err != nil {
			c.mu.Unlock()
			return nil, err
		}
		if c.session != nil {
			session := c.session
			c.mu.Unlock()
			return session, nil
		}
		if c.inFlight != nil {
			attempt := c.inFlight
			c.mu.Unlock()
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
		if err := ctx.Err(); err != nil {
			c.mu.Unlock()
			return nil, err
		}
		if c.connect == nil {
			c.mu.Unlock()
			return nil, ErrNoConnector
		}
		connectCtx, cancel := context.WithCancel(ctx)
		attempt := &connectAttempt{done: make(chan struct{}), cancel: cancel}
		c.inFlight = attempt
		c.mu.Unlock()

		session, err := c.connect(connectCtx)
		connectErr := connectCtx.Err()
		cancel()
		if err != nil && session != nil {
			_ = session.Close()
			session = nil
		}
		if err == nil {
			if connectErr != nil {
				ctxErr := connectErr
				if session != nil {
					_ = session.Close()
				}
				session = nil
				err = ctxErr
			} else if session == nil {
				err = ErrNilSession
			}
		}
		c.mu.Lock()
		if err == nil {
			c.session = session
		}
		attempt.err = err
		c.inFlight = nil
		close(attempt.done)
		c.mu.Unlock()
		if err != nil {
			return nil, err
		}
		return session, nil
	}
}
