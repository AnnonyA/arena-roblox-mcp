package agent

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestToolDispatcherNilReceiverReturnsError(t *testing.T) {
	var dispatcher *ToolDispatcher

	_, err := dispatcher.Dispatch(context.Background(), "script_read", json.RawMessage(`{}`))
	if !errors.Is(err, ErrNoToolDispatcher) {
		t.Fatalf("Dispatch() error = %v, want ErrNoToolDispatcher", err)
	}
}
