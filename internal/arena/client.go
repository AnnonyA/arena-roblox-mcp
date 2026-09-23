package arena

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const maxErrorBodyDrain = 64 * 1024

type ClientOptions struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewClient(opts ClientOptions) *Client {
	hc := opts.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: 60 * time.Second}
	}
	// Keep Arena credentials on the configured origin. In particular, do not
	// allow an API response to redirect an authenticated request to another
	// host. Clone caller-provided clients so enforcing this does not mutate
	// shared client configuration outside this package.
	client := *hc
	client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	baseURL := strings.TrimSpace(opts.BaseURL)
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{baseURL: baseURL, apiKey: opts.APIKey, http: &client}
}

func arenaStatusError(operation string, status int) error {
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return fmt.Errorf("Arena authentication failed. Check your API key.")
	}
	if status == http.StatusTooManyRequests {
		return fmt.Errorf("Arena rate limit exceeded. Try again later.")
	}
	return fmt.Errorf("Arena %s request failed with HTTP %d", operation, status)
}

func drainAndClose(body io.ReadCloser) {
	_, _ = io.CopyN(io.Discard, body, maxErrorBodyDrain)
	_ = body.Close()
}
