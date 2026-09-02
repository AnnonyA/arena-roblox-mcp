package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStreamChatAssemblesTextAndToolCallFragments(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret-key" {
			t.Fatalf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"Hel\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"lo\",\"tool_calls\":[{\"index\":0,\"id\":\"call_1\",\"type\":\"function\",\"function\":{\"name\":\"script_\",\"arguments\":\"{\\\"path\\\":\\\"A\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"function\":{\"name\":\"read\",\"arguments\":\"\\\"}\"}}]}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL, APIKey: "secret-key"})
	got, err := c.StreamChat(context.Background(), ChatRequest{
		Model:    "model-a",
		Messages: []Message{{Role: "user", Content: "inspect"}},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "Hello" {
		t.Fatalf("text = %q", got.Text)
	}
	if len(got.ToolCalls) != 1 {
		t.Fatalf("tool calls = %#v", got.ToolCalls)
	}
	call := got.ToolCalls[0]
	if call.ID != "call_1" || call.Type != "function" || call.Function.Name != "script_read" || call.Function.Arguments != `{"path":"A"}` {
		t.Fatalf("tool call = %#v", call)
	}
}

func TestStreamChatEmitsTextDeltas(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"one\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\" two\"}}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	var deltas []string
	_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, func(delta string) {
		deltas = append(deltas, delta)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(deltas) != 2 || deltas[0] != "one" || deltas[1] != " two" {
		t.Fatalf("deltas = %#v", deltas)
	}
}
