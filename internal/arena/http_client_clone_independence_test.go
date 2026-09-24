package arena

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClientCloneRemainsIndependentFromCallerMutation(t *testing.T) {
	t.Parallel()

	transport := &http.Transport{MaxIdleConns: 5}
	caller := &http.Client{
		Transport: transport,
		Timeout:   17 * time.Second,
	}
	arenaClient := NewClient(ClientOptions{
		BaseURL:    "https://example.invalid",
		HTTPClient: caller,
	})

	caller.Transport = http.DefaultTransport
	caller.Timeout = 3 * time.Second
	caller.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return nil }

	if arenaClient.http.Transport != transport {
		t.Fatal("Arena client transport changed after caller mutation")
	}
	if got, want := arenaClient.http.Timeout, 17*time.Second; got != want {
		t.Fatalf("Arena client timeout = %v after caller mutation, want %v", got, want)
	}
	if err := arenaClient.http.CheckRedirect(nil, nil); err != http.ErrUseLastResponse {
		t.Fatalf("Arena redirect policy error = %v after caller mutation, want http.ErrUseLastResponse", err)
	}
}
