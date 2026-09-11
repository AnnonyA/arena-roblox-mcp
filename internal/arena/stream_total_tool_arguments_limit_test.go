package arena

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatRejectsOversizedTotalToolCallArguments(t *testing.T) {
	t.Parallel()

	payload := strings.Repeat("a", 600*1024)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		for i := 0; i < 2; i++ {
			arguments := `{"data":"` + payload + `"}`
			_, _ = fmt.Fprintf(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":%d,\"id\":\"call_%d\",\"type\":\"function\",\"function\":{\"name\":\"script_write\",\"arguments\":%q}}]}}]}\n\n", i, i, arguments)
		}
		_, _ = fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := c.StreamChat(context.Background(), ChatRequest{Model: "model-a"}, nil)
	if err == nil {
		t.Fatal("StreamChat error = nil, want total tool-call arguments limit error")
	}
	if !strings.Contains(err.Error(), "total tool call arguments exceed") {
		t.Fatalf("StreamChat error = %q, want total tool-call arguments limit error", err)
	}
}
