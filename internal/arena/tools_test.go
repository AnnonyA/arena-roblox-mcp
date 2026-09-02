package arena

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStreamChatSendsToolDefinitions(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var got struct {
			Stream bool             `json:"stream"`
			Tools  []ToolDefinition `json:"tools"`
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		if !got.Stream {
			t.Fatal("stream = false, want true")
		}
		if len(got.Tools) != 1 {
			t.Fatalf("tools = %#v", got.Tools)
		}
		tool := got.Tools[0]
		if tool.Type != "function" || tool.Function.Name != "script_read" {
			t.Fatalf("tool = %#v", tool)
		}
		if string(tool.Function.Parameters) != `{"type":"object","properties":{"path":{"type":"string"}}}` {
			t.Fatalf("parameters = %s", tool.Function.Parameters)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.StreamChat(context.Background(), ChatRequest{
		Model: "model-a",
		Tools: []ToolDefinition{{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "script_read",
				Description: "Read a Luau script",
				Parameters:  json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"}}}`),
			},
		}},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
