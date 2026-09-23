package arena

import (
	"context"
	"strings"
	"testing"
)

func TestStreamChatRejectsInvalidClientState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *Client
		ctx     context.Context
		wantErr string
	}{
		{name: "nil client", client: nil, ctx: context.Background(), wantErr: "stream Arena chat: client is nil"},
		{name: "nil context", client: NewClient(ClientOptions{}), ctx: nil, wantErr: "stream Arena chat: context is nil"},
		{name: "userinfo", client: NewClient(ClientOptions{BaseURL: "https://user:password@api.arena.ai"}), ctx: context.Background(), wantErr: "stream Arena chat: base URL must not contain userinfo"},
		{name: "query", client: NewClient(ClientOptions{BaseURL: "https://api.arena.ai?tenant=test"}), ctx: context.Background(), wantErr: "stream Arena chat: base URL must not contain query or fragment"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.client.StreamChat(tt.ctx, ChatRequest{Model: "test"}, nil)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want substring %q", err, tt.wantErr)
			}
		})
	}
}
