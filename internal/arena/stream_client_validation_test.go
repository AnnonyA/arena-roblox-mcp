package arena

import (
	"context"
	"net/http"
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
		{name: "nil client", client: nil, ctx: context.Background(), want: "stream Arena chat: client is nil"},
		{name: "nil context", client: NewClient(ClientOptions{BaseURL: "https://example.com"}), ctx: nil, want: "stream Arena chat: context is nil"},
		{name: "nil HTTP client", client: &Client{baseURL: "https://example.com"}, ctx: context.Background(), want: "stream Arena chat: HTTP client is nil"},
		{name: "empty base URL", client: &Client{http: http.DefaultClient}, ctx: context.Background(), want: "stream Arena chat: base URL is empty"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.client.StreamChat(tt.ctx, ChatRequest{Model: "test"}, nil)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if got := err.Error(); got != tt.want {
				t.Fatalf("error = %q, want %q", got, tt.want)
			}
		})
	}
}
