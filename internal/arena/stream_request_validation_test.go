package arena

import (
	"context"
	"net/http"
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
		{name: "nil client", client: nil, ctx: context.Background(), want: "stream Arena chat: client is nil"},
		{name: "nil context", client: NewClient(ClientOptions{BaseURL: "https://example.com"}), ctx: nil, want: "stream Arena chat: context is nil"},
		{name: "nil HTTP client", client: &Client{baseURL: "https://example.com"}, ctx: context.Background(), want: "stream Arena chat: HTTP client is nil"},
		{name: "empty base URL", client: NewClient(ClientOptions{}), ctx: context.Background(), want: "stream Arena chat: base URL is empty"},
		{name: "invalid scheme", client: NewClient(ClientOptions{BaseURL: "ftp://example.com"}), ctx: context.Background(), want: "stream Arena chat: base URL must use http or https"},
		{name: "userinfo", client: NewClient(ClientOptions{BaseURL: "https://user:pass@example.com"}), ctx: context.Background(), want: "stream Arena chat: base URL must not contain userinfo"},
		{name: "query", client: NewClient(ClientOptions{BaseURL: "https://example.com?target=other"}), ctx: context.Background(), want: "stream Arena chat: base URL must not contain query or fragment"},
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

func TestStreamChatValidationDoesNotPerformRequest(t *testing.T) {
	t.Parallel()

	called := false
	client := NewClient(ClientOptions{
		BaseURL: "https://user:pass@example.com",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			called = true
			return nil, nil
		})},
	})

	_, err := client.StreamChat(context.Background(), ChatRequest{Model: "test"}, nil)
	if err == nil || !strings.Contains(err.Error(), "userinfo") {
		t.Fatalf("error = %v, want userinfo validation error", err)
	}
	if called {
		t.Fatal("HTTP transport was called for invalid base URL")
	}
}
