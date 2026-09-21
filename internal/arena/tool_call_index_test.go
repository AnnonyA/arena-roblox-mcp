package arena

import (
	"strings"
	"testing"
)

func TestToolCallUnmarshalRejectsNegativeIndex(t *testing.T) {
	data := []byte(`{"index":-1,"id":"call_1","type":"function","function":{"name":"read_script","arguments":"{}"}}`)

	var call ToolCall
	err := call.UnmarshalJSON(data)
	if err == nil {
		t.Fatal("expected negative tool call index to be rejected")
	}
	if !strings.Contains(err.Error(), "negative tool call index") {
		t.Fatalf("error = %q, want negative index validation error", err)
	}
}

func TestToolCallUnmarshalAcceptsZeroIndex(t *testing.T) {
	data := []byte(`{"index":0,"id":"call_1","type":"function","function":{"name":"read_script","arguments":"{}"}}`)

	var call ToolCall
	if err := call.UnmarshalJSON(data); err != nil {
		t.Fatalf("expected zero tool call index to be accepted: %v", err)
	}
	if call.Index != 0 {
		t.Fatalf("index = %d, want 0", call.Index)
	}
}
