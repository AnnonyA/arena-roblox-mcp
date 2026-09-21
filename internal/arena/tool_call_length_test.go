package arena

import (
	"strings"
	"testing"
)

func TestToolCallUnmarshalRejectsOversizedID(t *testing.T) {
	data := []byte(`{"id":"` + strings.Repeat("a", maxToolCallIdentifierBytes+1) + `","type":"function","function":{"name":"read_script","arguments":"{}"}}`)

	var call ToolCall
	err := call.UnmarshalJSON(data)
	if err == nil {
		t.Fatal("expected oversized tool call id to be rejected")
	}
	if !strings.Contains(err.Error(), "id exceeds") {
		t.Fatalf("error = %q, want id size validation error", err)
	}
}

func TestToolCallUnmarshalRejectsOversizedName(t *testing.T) {
	data := []byte(`{"id":"call_1","type":"function","function":{"name":"` + strings.Repeat("a", maxToolCallIdentifierBytes+1) + `","arguments":"{}"}}`)

	var call ToolCall
	err := call.UnmarshalJSON(data)
	if err == nil {
		t.Fatal("expected oversized tool call name to be rejected")
	}
	if !strings.Contains(err.Error(), "name exceeds") {
		t.Fatalf("error = %q, want name size validation error", err)
	}
}

func TestToolCallUnmarshalRejectsOversizedArguments(t *testing.T) {
	arguments := strings.Repeat("a", maxToolCallArgumentsBytes+1)
	data := []byte(`{"id":"call_1","type":"function","function":{"name":"write_script","arguments":"` + arguments + `"}}`)

	var call ToolCall
	err := call.UnmarshalJSON(data)
	if err == nil {
		t.Fatal("expected oversized tool call arguments to be rejected")
	}
	if !strings.Contains(err.Error(), "arguments exceed") {
		t.Fatalf("error = %q, want arguments size validation error", err)
	}
}

func TestToolCallUnmarshalAcceptsFieldsAtSizeLimits(t *testing.T) {
	id := strings.Repeat("a", maxToolCallIdentifierBytes)
	name := strings.Repeat("b", maxToolCallIdentifierBytes)
	arguments := strings.Repeat("c", maxToolCallArgumentsBytes)
	data := []byte(`{"id":"` + id + `","type":"function","function":{"name":"` + name + `","arguments":"` + arguments + `"}}`)

	var call ToolCall
	if err := call.UnmarshalJSON(data); err != nil {
		t.Fatalf("expected tool call fields at size limits to be accepted: %v", err)
	}
	if call.ID != id || call.Function.Name != name || call.Function.Arguments != arguments {
		t.Fatal("unexpected decoded tool call fields")
	}
}
