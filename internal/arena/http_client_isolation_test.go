package arena

import (
	"net/http"
	"testing"
)

func TestNewClientDoesNotMutateCallerHTTPClient(t *testing.T) {
	t.Parallel()

	redirectCalls := 0
	callerPolicy := func(_ *http.Request, _ []*http.Request) error {
		redirectCalls++
		return nil
	}
	caller := &http.Client{CheckRedirect: callerPolicy}

	arenaClient := NewClient(ClientOptions{
		BaseURL:    "https://example.invalid",
		APIKey:     "test-key",
		HTTPClient: caller,
	})

	if caller.CheckRedirect == nil {
		t.Fatal("NewClient cleared caller redirect policy")
	}
	if err := caller.CheckRedirect(nil, nil); err != nil {
		t.Fatalf("caller redirect policy returned error after NewClient: %v", err)
	}
	if redirectCalls != 1 {
		t.Fatalf("caller redirect policy calls = %d, want 1", redirectCalls)
	}
	if arenaClient.http == caller {
		t.Fatal("Arena client reused caller *http.Client instead of cloning it")
	}
	if err := arenaClient.http.CheckRedirect(nil, nil); err != http.ErrUseLastResponse {
		t.Fatalf("Arena redirect policy error = %v, want http.ErrUseLastResponse", err)
	}
	if redirectCalls != 1 {
		t.Fatalf("Arena redirect policy invoked caller policy; calls = %d, want 1", redirectCalls)
	}
}
