package arena

import (
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"
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

func TestNewClientClonePreservesCallerHTTPSettings(t *testing.T) {
	t.Parallel()

	transport := &http.Transport{MaxIdleConns: 7}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("create cookie jar: %v", err)
	}
	caller := &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   23 * time.Second,
	}

	arenaClient := NewClient(ClientOptions{
		BaseURL:    "https://example.invalid",
		HTTPClient: caller,
	})

	if arenaClient.http.Transport != transport {
		t.Fatal("Arena client did not preserve caller transport")
	}
	if arenaClient.http.Jar != jar {
		t.Fatal("Arena client did not preserve caller cookie jar")
	}
	if arenaClient.http.Timeout != caller.Timeout {
		t.Fatalf("Arena client timeout = %v, want %v", arenaClient.http.Timeout, caller.Timeout)
	}
	if caller.CheckRedirect != nil {
		t.Fatal("NewClient unexpectedly installed redirect policy on caller client")
	}
	if arenaClient.http.CheckRedirect == nil {
		t.Fatal("Arena client did not install redirect protection on clone")
	}
}
