package arena

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type trackingBody struct {
	reader    *strings.Reader
	eof       bool
	closed    bool
	bytesRead int
}

func newTrackingBody(content string) *trackingBody {
	return &trackingBody{reader: strings.NewReader(content)}
}

func (b *trackingBody) Read(p []byte) (int, error) {
	n, err := b.reader.Read(p)
	b.bytesRead += n
	if err == io.EOF {
		b.eof = true
	}
	return n, err
}

func (b *trackingBody) Close() error {
	b.closed = true
	return nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestListModelsDrainsTerminalErrorBody(t *testing.T) {
	body := newTrackingBody("denied")
	client := NewClient(ClientOptions{
		BaseURL: "https://arena.invalid",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusUnauthorized, Body: body, Header: make(http.Header)}, nil
		})},
	})

	_, err := client.ListModels(context.Background())
	if err == nil {
		t.Fatal("ListModels error = nil, want authentication error")
	}
	if !body.eof {
		t.Fatal("ListModels did not drain terminal error body to EOF")
	}
	if !body.closed {
		t.Fatal("ListModels did not close terminal error body")
	}
}

func TestListModelsBoundsTerminalErrorBodyDrain(t *testing.T) {
	const maxExpectedDrain = 64 * 1024
	body := newTrackingBody(strings.Repeat("x", maxExpectedDrain*4))
	client := NewClient(ClientOptions{
		BaseURL: "https://arena.invalid",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusUnauthorized, Body: body, Header: make(http.Header)}, nil
		})},
	})

	_, err := client.ListModels(context.Background())
	if err == nil {
		t.Fatal("ListModels error = nil, want authentication error")
	}
	if body.bytesRead > maxExpectedDrain {
		t.Fatalf("ListModels drained %d bytes from oversized error body, want at most %d", body.bytesRead, maxExpectedDrain)
	}
	if !body.closed {
		t.Fatal("ListModels did not close oversized terminal error body")
	}
}

func TestStreamChatDrainsTerminalErrorBody(t *testing.T) {
	body := newTrackingBody("denied")
	client := NewClient(ClientOptions{
		BaseURL: "https://arena.invalid",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusUnauthorized, Body: body, Header: make(http.Header)}, nil
		})},
	})

	_, err := client.StreamChat(context.Background(), ChatRequest{Model: "test-model"}, nil)
	if err == nil {
		t.Fatal("StreamChat error = nil, want authentication error")
	}
	if !body.eof {
		t.Fatal("StreamChat did not drain terminal error body to EOF")
	}
	if !body.closed {
		t.Fatal("StreamChat did not close terminal error body")
	}
}
