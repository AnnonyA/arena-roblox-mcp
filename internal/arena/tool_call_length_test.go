package arena

import (
	"strings"
	"testing"
)

func TestToolCallUnmarshalAcceptsPartialStreamIdentifiers(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{name: "missing id fragment", data: `{"id":"","type":"function","function":{"name":"read_script","arguments":""}}`},
		{name: "missing name fragment", data: `{"id":"call_1","type":"function","function":{"name":"","arguments":""}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var call ToolCall
			if err := call.UnmarshalJSON([]byte(tt.data)); err != nil {
				t.Fatalf("partial streaming tool call should decode before final assembly validation: %v", err)
			}
		})
	}
}

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

func TestToolCallUnmarshalRejectsOversizedType(t *testing.T) {
	data := []byte(`{"id":"call_1","type":"` + strings.Repeat("a", maxToolCallIdentifierBytes+1) + `","function":{"name":"read_script","arguments":"{}"}}`)

	var call ToolCall
	err := call.UnmarshalJSON(data)
	if err == nil {
		t.Fatal("expected oversized tool call type to be rejected")
	}
	if !strings.Contains(err.Error(), "type exceeds") {
		t.Fatalf("error = %q, want type size validation error", err)
	}
}

func TestToolCallUnmarshalRejectsMalformedType(t *testing.T) {
	tests := []struct {
		name     string
		callType string
		want     string
	}{
		{name: "surrounding whitespace", callType: " function", want: "surrounding whitespace"},
		{name: "control character", callType: "func\\ntion", want: "control character"},
		{name: "line separator", callType: "func\\u2028tion", want: "line separator"},
		{name: "format character", callType: "func\\u200Btion", want: "format character"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(`{"id":"call_1","type":"` + tt.callType + `","function":{"name":"read_script","arguments":"{}"}}`)
			var call ToolCall
			err := call.UnmarshalJSON(data)
			if err == nil {
				t.Fatal("expected malformed tool call type to be rejected")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want %q", err, tt.want)
			}
		})
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
	callType := strings.Repeat("t", maxToolCallIdentifierBytes)
	name := strings.Repeat("b", maxToolCallIdentifierBytes)
	arguments := strings.Repeat("c", maxToolCallArgumentsBytes)
	data := []byte(`{"id":"` + id + `","type":"` + callType + `","function":{"name":"` + name + `","arguments":"` + arguments + `"}}`)

	var call ToolCall
	if err := call.UnmarshalJSON(data); err != nil {
		t.Fatalf("expected tool call fields at size limits to be accepted: %v", err)
	}
	if call.ID != id || call.Type != callType || call.Function.Name != name || call.Function.Arguments != arguments {
		t.Fatal("unexpected decoded tool call fields")
	}
}
