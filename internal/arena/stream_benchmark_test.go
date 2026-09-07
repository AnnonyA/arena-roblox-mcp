package arena

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkStreamParser(b *testing.B) {
	payload := []byte("data: {\"choices\":[{\"delta\":{\"content\":\"chunk\"}}]}\n\ndata: {\"choices\":[{\"delta\":{\"content\":\"chunk\"}}]}\n\ndata: [DONE]\n\n")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	client := NewClient(ClientOptions{BaseURL: srv.URL})
	request := ChatRequest{Model: "arena/benchmark", Messages: []Message{{Role: "user", Content: "benchmark"}}}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := client.StreamChat(ctx, request, nil)
		if err != nil {
			b.Fatal(err)
		}
		if result.Text != "chunkchunk" {
			b.Fatalf("text = %q", result.Text)
		}
	}
}
