package arena

import (
	"context"
	"testing"
)

func TestStreamChatRejectsNilClient(t *testing.T) {
	t.Parallel()

	var c *Client
	_, err := c.StreamChat(context.Background(), ChatRequest{}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	const want = "stream Arena chat: client is nil"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err, want)
	}
}
