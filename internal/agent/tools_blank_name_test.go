package agent

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestToolDispatcherRejectsBlankToolNameWithoutCallingBackend(t *testing.T) {
	for _, name := range []string{"", "   ", "\t\n"} {
		caller := &recordingToolCaller{}
		dispatcher := NewToolDispatcher([]string{name}, caller)

		_, err := dispatcher.Dispatch(context.Background(), name, json.RawMessage(`{}`))

		if !errors.Is(err, ErrUnknownTool) {
			t.Fatalf("Dispatch(%q) error = %v, want ErrUnknownTool", name, err)
		}
		if caller.calls != 0 {
			t.Fatalf("Dispatch(%q) backend calls = %d, want 0", name, caller.calls)
		}
	}
}
