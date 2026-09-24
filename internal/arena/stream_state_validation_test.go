package arena

import (
	"context"
	"strings"
	"testing"
)

func TestStreamChatRejectsInvalidClientState(t *testing.T) {
	t.Parallel()

	t.Run("nil client", func(t *testing.T) {
		var client *Client
		_, err := client.StreamChat(context.Background(), ChatRequest{}, nil)
		if err == nil || !strings.Contains(err.Error(), "client is nil") {
			t.Fatalf("nil Client.StreamChat error = %v, want client is nil", err)
		}
	})

	t.Run("nil context", func(t *testing.T) {
		client := NewClient(ClientOptions{BaseURL: "https://example.invalid"})
		_, err := client.StreamChat(nil, ChatRequest{}, nil)
		if err == nil || !strings.Contains(err.Error(), "context is nil") {
			t.Fatalf("nil context StreamChat error = %v, want context is nil", err)
		}
	})

	t.Run("nil HTTP client", func(t *testing.T) {
		client := &Client{baseURL: "https://example.invalid"}
		_, err := client.StreamChat(context.Background(), ChatRequest{}, nil)
		if err == nil || !strings.Contains(err.Error(), "HTTP client is nil") {
			t.Fatalf("nil HTTP client StreamChat error = %v, want HTTP client is nil", err)
		}
	})
}
