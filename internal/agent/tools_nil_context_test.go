package agent

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestToolDispatcherRejectsNilContextWithoutCallingBackend(t *testing.T) {
	caller := &recordingToolCaller{}
	dispatcher := NewToolDispatcher([]string{"read_script"}, caller)

	_, err := dispatcher.Dispatch(nil, "read_script", json.RawMessage(`{}`))

	if !errors.Is(err, ErrNoContext) {
		t.Fatalf("Dispatch() error = %v, want ErrNoContext", err)
	}
	if caller.calls != 0 {
		t.Fatalf("backend calls = %d, want 0", caller.calls)
	}
}
