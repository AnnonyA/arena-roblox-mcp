package arena

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func FuzzToolCallValidation(f *testing.F) {
	for _, seed := range []struct {
		id       string
		callType string
		name     string
	}{
		{id: "call-1", callType: "function", name: "inspect"},
		{id: " call-1", callType: "function", name: "inspect"},
		{id: "call\n1", callType: "function", name: "inspect"},
		{id: "call\u200b1", callType: "function", name: "inspect"},
		{id: "call\u20281", callType: "function", name: "inspect"},
		{id: "call-1", callType: " function", name: "inspect"},
		{id: "call-1", callType: "func\u2060tion", name: "inspect"},
		{id: "call-1", callType: "function\u2029", name: "inspect"},
		{id: "call-1", callType: "function", name: "ins\u2060pect"},
		{id: "call-1", callType: "function", name: "inspect\u2029"},
		{id: "呼び出し-1", callType: "関数", name: "検査"},
	} {
		f.Add(seed.id, seed.callType, seed.name)
	}

	f.Fuzz(func(t *testing.T, id, callType, name string) {
		if !utf8.ValidString(id) || !utf8.ValidString(callType) || !utf8.ValidString(name) {
			t.Skip()
		}

		payload, err := json.Marshal(map[string]any{
			"id":   id,
			"type": callType,
			"function": map[string]any{
				"name":      name,
				"arguments": "{}",
			},
		})
		if err != nil {
			t.Fatalf("json.Marshal: %v", err)
		}

		var call ToolCall
		if err := json.Unmarshal(payload, &call); err != nil {
			return
		}
		if call.ID != id || call.Type != callType || call.Function.Name != name {
			t.Fatalf("json.Unmarshal changed tool call fields: got id=%q type=%q name=%q, want id=%q type=%q name=%q", call.ID, call.Type, call.Function.Name, id, callType, name)
		}

		for field, value := range map[string]string{"id": id, "type": callType, "name": name} {
			if strings.TrimSpace(value) != value {
				t.Fatalf("json.Unmarshal accepted surrounding whitespace in tool call %s %q", field, value)
			}
			for _, r := range value {
				if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || r == '\u2028' || r == '\u2029' {
					t.Fatalf("json.Unmarshal accepted unsafe rune U+%04X in tool call %s %q", r, field, value)
				}
			}
		}
	})
}
