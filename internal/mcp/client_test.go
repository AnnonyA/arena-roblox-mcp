package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"
)

type fakeSession struct {
	mu         sync.Mutex
	listCalls  int
	closeCalls int
}

func (s *fakeSession) ListTools(context.Context) ([]Tool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listCalls++
	return []Tool{{Name: "read_script"}}, nil
}

func (s *fakeSession) CallTool(context.Context, string, json.RawMessage) (ToolResult, error) {
	return ToolResult{}, nil
}

func (s *fakeSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeCalls++
	return nil
}

func TestClientReusesPersistentSessionAndToolDiscovery(t *testing.T) {
	session := &fakeSession{}
	connectCalls := 0
	client := NewClient(func(context.Context) (Session, error) {
		connectCalls++
		return session, nil
	})

	first, err := client.Tools(context.Background())
	if err != nil {
		t.Fatalf("first Tools: %v", err)
	}
	second, err := client.Tools(context.Background())
	if err != nil {
		t.Fatalf("second Tools: %v", err)
	}

	if connectCalls != 1 {
		t.Fatalf("connect calls = %d, want 1", connectCalls)
	}
	if session.listCalls != 1 {
		t.Fatalf("ListTools calls = %d, want 1", session.listCalls)
	}
	if len(first) != 1 || len(second) != 1 || first[0].Name != "read_script" || second[0].Name != "read_script" {
		t.Fatalf("unexpected tools: first=%v second=%v", first, second)
	}
}

func TestClientCloseClosesSession(t *testing.T) {
	session := &fakeSession{}
	client := NewClient(func(context.Context) (Session, error) {
		return session, nil
	})

	if _, err := client.Tools(context.Background()); err != nil {
		t.Fatalf("Tools before Close: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if session.closeCalls != 1 {
		t.Fatalf("session Close calls = %d, want 1", session.closeCalls)
	}
}

func TestClientWaiterCanCancelWhileConnectionIsInFlight(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	firstDone := make(chan error, 1)
	session := &fakeSession{}

	client := NewClient(func(context.Context) (Session, error) {
		close(started)
		<-release
		return session, nil
	})

	go func() {
		_, err := client.CallTool(context.Background(), "read_script", nil)
		firstDone <- err
	}()

	<-started

	ctx, cancel := context.WithCancel(context.Background())
	secondDone := make(chan error, 1)
	go func() {
		_, err := client.CallTool(ctx, "read_script", nil)
		secondDone <- err
	}()
	cancel()

	select {
	case err := <-secondDone:
		if !errors.Is(err, context.Canceled) {
			close(release)
			<-firstDone
			t.Fatalf("second CallTool error = %v, want context.Canceled", err)
		}
	case <-time.After(100 * time.Millisecond):
		close(release)
		<-firstDone
		t.Fatal("second CallTool ignored cancellation while connection was in progress")
	}

	close(release)
	if err := <-firstDone; err != nil {
		t.Fatalf("first CallTool: %v", err)
	}
}

func TestClientRejectsMissingConnector(t *testing.T) {
	client := NewClient(nil)

	_, err := client.Tools(context.Background())
	if !errors.Is(err, ErrNoConnector) {
		t.Fatalf("Tools error = %v, want ErrNoConnector", err)
	}
}

func TestClientRejectsNilSessionFromConnector(t *testing.T) {
	client := NewClient(func(context.Context) (Session, error) {
		return nil, nil
	})

	_, err := client.Tools(context.Background())
	if !errors.Is(err, ErrNilSession) {
		t.Fatalf("Tools error = %v, want ErrNilSession", err)
	}
}
