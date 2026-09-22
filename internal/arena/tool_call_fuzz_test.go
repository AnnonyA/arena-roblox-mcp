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
		index     int
		id        string
		callType  string
		name      string
		arguments string
	}{
		{index: 0, id: "call-1", callType: "function", name: "inspect", arguments: "{}"},
		{index: 1, id: "call-1", callType: "function", name: "inspect", arguments: "{}"},
		{index: -1, id: "call-1", callType: "function", name: "inspect", arguments: "{}"},
		{index: 0, id: " call-1", callType: "function", name: "inspect", arguments: "{}"},
		{index: 0, id: "call\n1", callType: "function", name: "inspect", arguments: "{}"},
		{index: 0, id: "call\u200b1", callType: "function", name: "inspect", arguments: "{}"},
		{index: 0, id: "call\u20281", callType: "function", name: "inspect", arguments: "{}"},
		{index: 0, id: "call-1", callType: " function", name: "inspect", arguments: "{}"},
		{index: 0, id: "call-1", callType: "func\u2060tion", name: "inspect", arguments: "{}"},
		{index: 0, id: "call-1", callType: "function\u2029", name: "inspect", arguments: "{}"},
		{index: 0, id: "call-1", callType: "function", name: "ins\u2060pect", arguments: "{}"},
		{index: 0, id: "call-1", callType: "function", name: "inspect\u2029", arguments: "{}"},
		{index: 0, id: "呼び出し-1", callType: "関数", name: "検査", arguments: `{"対象":"Workspace"}`},
		{index: 0, id: "call-1", callType: "function", name: "inspect", arguments: ""},
		{index: 0, id: "call-1", callType: "function", name: "inspect", arguments: `{"path":"Workspace.Part"}`},
	} {
		f.Add(seed.index, seed.id, seed.callType, seed.name, seed.arguments)
	}

	f.Fuzz(func(t *testing.T, index int, id, callType, name, arguments string) {
		if !utf8.ValidString(id) || !utf8.ValidString(callType) || !utf8.ValidString(name) || !utf8.ValidString(arguments) {
			t.Skip()
		}

		payload, err := json.Marshal(map[string]any{
			"index": index,
			"id":    id,
			"type":  callType,
			"function": map[string]any{
				"name":      name,
				"arguments": arguments,
			},
		})
		if err != nil {
			t.Fatalf("json.Marshal: %v", err)
		}

		var call ToolCall
		if err := json.Unmarshal(payload, &call); err != nil {
			return
		}
		if call.Index != index || call.ID != id || call.Type != callType || call.Function.Name != name || call.Function.Arguments != arguments {
			t.Fatalf("json.Unmarshal changed tool call fields: got index=%d id=%q type=%q name=%q arguments=%q, want index=%d id=%q type=%q name=%q arguments=%q", call.Index, call.ID, call.Type, call.Function.Name, call.Function.Arguments, index, id, callType, name, arguments)
		}
		if index < 0 {
			t.Fatalf("json.Unmarshal accepted negative tool call index %d", index)
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
