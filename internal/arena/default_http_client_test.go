package arena

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClientDefaultHTTPClientHasBoundedTimeoutAndRedirectProtection(t *testing.T) {
	t.Parallel()

	client := NewClient(ClientOptions{BaseURL: "https://example.invalid"})
	if client.http == nil {
		t.Fatal("NewClient did not create a default HTTP client")
	}
	if client.http.Timeout != 60*time.Second {
		t.Fatalf("default HTTP timeout = %v, want %v", client.http.Timeout, 60*time.Second)
	}
	if client.http.CheckRedirect == nil {
		t.Fatal("default HTTP client has no redirect protection")
	}
	if err := client.http.CheckRedirect(nil, nil); err != http.ErrUseLastResponse {
		t.Fatalf("default redirect policy error = %v, want http.ErrUseLastResponse", err)
	}
}
