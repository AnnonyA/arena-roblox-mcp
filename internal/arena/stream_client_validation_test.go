package arena

import (
	"context"
	"strings"
	"testing"
)

func TestStreamChatRejectsInvalidClientState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		client *Client
		ctx    context.Context
		want   string
	}{
		{name: "nil client", client: nil, ctx: context.Background(), want: "client is nil"},
		{name: "nil context", client: NewClient(ClientOptions{BaseURL: "https://api.arena.ai"}), ctx: nil, want: "context is nil"},
		{name: "nil HTTP client", client: &Client{baseURL: "https://api.arena.ai"}, ctx: context.Background(), want: "HTTP client is nil"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.client.StreamChat(tt.ctx, ChatRequest{Model: "model-a"}, nil)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want substring %q", err, tt.want)
			}
		})
	}
}
