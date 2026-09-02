package arena

import (
	"net/http"
	"strings"
	"time"
)

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
