package mcp

import (
	"context"
	"sync"
	"testing"
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

func TestClientCloseClosesSessionAndForcesReconnect(t *testing.T) {
	first := &fakeSession{}
	second := &fakeSession{}
	connectCalls := 0
	client := NewClient(func(context.Context) (Session, error) {
		connectCalls++
		if connectCalls == 1 {
			return first, nil
		}
		return second, nil
	})

	if _, err := client.Tools(context.Background()); err != nil {
		t.Fatalf("Tools before Close: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, err := client.Tools(context.Background()); err != nil {
		t.Fatalf("Tools after Close: %v", err)
	}

	if first.closeCalls != 1 {
		t.Fatalf("first session Close calls = %d, want 1", first.closeCalls)
	}
	if connectCalls != 2 {
		t.Fatalf("connect calls = %d, want 2", connectCalls)
	}
	if second.listCalls != 1 {
		t.Fatalf("second session ListTools calls = %d, want 1", second.listCalls)
	}
}
