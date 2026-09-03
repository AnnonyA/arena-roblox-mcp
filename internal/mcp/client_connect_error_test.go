package mcp

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientSharesConnectionErrorWithConcurrentWaiters(t *testing.T) {
	wantErr := errors.New("studio unavailable")
	started := make(chan struct{})
	release := make(chan struct{})
	var connectCalls atomic.Int32

	client := NewClient(func(context.Context) (Session, error) {
		if connectCalls.Add(1) == 1 {
			close(started)
		}
		<-release
		return nil, wantErr
	})

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := client.CallTool(context.Background(), "read_script", nil)
		errs <- err
	}()
	<-started
	go func() {
		defer wg.Done()
		_, err := client.CallTool(context.Background(), "read_script", nil)
		errs <- err
	}()

	time.Sleep(20 * time.Millisecond)
	close(release)
	wg.Wait()
	close(errs)

	for err := range errs {
		if !errors.Is(err, wantErr) {
			t.Fatalf("CallTool error = %v, want %v", err, wantErr)
		}
	}
	if got := connectCalls.Load(); got != 1 {
		t.Fatalf("connect calls = %d, want 1", got)
	}
}
