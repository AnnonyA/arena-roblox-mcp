package arena

import (
	"context"
	"strings"
	"testing"
)

func TestListModelsRejectsInvalidClientState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *Client
		ctx     context.Context
		wantErr string
	}{
		{name: "nil client", client: nil, ctx: context.Background(), wantErr: "client is nil"},
		{name: "nil context", client: NewClient(ClientOptions{BaseURL: "https://arena.example"}), ctx: nil, wantErr: "context is nil"},
		{name: "nil HTTP client", client: &Client{baseURL: "https://arena.example"}, ctx: context.Background(), wantErr: "HTTP client is nil"},
		{name: "empty base URL", client: &Client{http: NewClient(ClientOptions{BaseURL: "https://arena.example"}).http}, ctx: context.Background(), wantErr: "base URL is empty"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.client.ListModels(tt.ctx)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("ListModels() error = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}
