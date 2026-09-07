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
	return &Client{baseURL: strings.TrimRight(opts.BaseURL, "/"), apiKey: opts.APIKey, http: hc}
}

func arenaStatusError(operation string, status int) error {
	if status == http.StatusUnauthorized {
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
