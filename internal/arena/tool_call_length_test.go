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

func TestToolCallUnmarshalAcceptsIdentifiersAtSizeLimit(t *testing.T) {
	id := strings.Repeat("a", maxToolCallIdentifierBytes)
	name := strings.Repeat("b", maxToolCallIdentifierBytes)
	data := []byte(`{"id":"` + id + `","type":"function","function":{"name":"` + name + `","arguments":"{}"}}`)

	var call ToolCall
	if err := call.UnmarshalJSON(data); err != nil {
		t.Fatalf("expected identifiers at size limit to be accepted: %v", err)
	}
	if call.ID != id || call.Function.Name != name {
		t.Fatalf("unexpected decoded identifiers: id=%q name=%q", call.ID, call.Function.Name)
	}
}
