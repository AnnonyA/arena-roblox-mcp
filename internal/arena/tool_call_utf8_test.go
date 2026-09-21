package arena

import (
	"strings"
	"testing"
)

func TestToolCallUnmarshalRejectsInvalidUTF8(t *testing.T) {
	data := append([]byte(`{"id":"call_`), 0xff)
	data = append(data, []byte(`","type":"function","function":{"name":"read_script","arguments":"{}"}}`)...)

	var call ToolCall
	err := call.UnmarshalJSON(data)
	if err == nil {
		t.Fatal("expected invalid UTF-8 tool call payload to be rejected")
	}
	if !strings.Contains(err.Error(), "invalid UTF-8") {
		t.Fatalf("error = %q, want invalid UTF-8 error", err)
	}
}
