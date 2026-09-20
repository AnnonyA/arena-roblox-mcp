package arena

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToolCallRejectsEmptyIdentityFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		fnName  string
		wantErr string
	}{
		{name: "empty id", id: "", fnName: "inspect", wantErr: "tool call id is empty"},
		{name: "empty name", id: "call-1", fnName: "", wantErr: "tool call name is empty"},
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
