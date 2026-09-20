package arena

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToolCallRejectsUnicodeFormatCharacters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		fnName  string
		wantErr string
	}{
		{name: "zero width space in id", id: "call\u200b1", fnName: "inspect", wantErr: "tool call id contains format character"},
		{name: "word joiner in id", id: "call\u20601", fnName: "inspect", wantErr: "tool call id contains format character"},
		{name: "zero width space in name", id: "call-1", fnName: "ins\u200bpect", wantErr: "tool call name contains format character"},
		{name: "word joiner in name", id: "call-1", fnName: "ins\u2060pect", wantErr: "tool call name contains format character"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			payload, err := json.Marshal(map[string]any{
				"id":   tt.id,
				"type": "function",
				"function": map[string]any{
					"name":      tt.fnName,
					"arguments": "{}",
				},
			})
			if err != nil {
				t.Fatalf("json.Marshal: %v", err)
			}

			var call ToolCall
			err = json.Unmarshal(payload, &call)
			if err == nil {
				t.Fatalf("json.Unmarshal error = nil, want %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("json.Unmarshal error = %q, want %q", err, tt.wantErr)
			}
		})
	}
}
