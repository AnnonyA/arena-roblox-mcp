package arena

import (
	"encoding/json"
	"testing"
)

func TestMessagePreservesToolCallFields(t *testing.T) {
	const raw = `{"role":"assistant","tool_calls":[{"id":"call-1","type":"function","function":{"name":"script_read","arguments":"{\"path\":\"ServerScriptService.Main\"}"}}],"tool_call_id":"call-1"}`

	var message Message
	if err := json.Unmarshal([]byte(raw), &message); err != nil {
		t.Fatalf("unmarshal message: %v", err)
	}

	encoded, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}

	var got map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("decode round-trip message: %v", err)
	}
	if _, ok := got["tool_calls"]; !ok {
		t.Fatalf("round-trip message dropped tool_calls: %s", encoded)
	}
	if _, ok := got["tool_call_id"]; !ok {
		t.Fatalf("round-trip message dropped tool_call_id: %s", encoded)
	}
}
