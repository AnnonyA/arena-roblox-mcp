package arena

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

type modelsPreflightTransport struct {
	called bool
}

func (t *modelsPreflightTransport) RoundTrip(*http.Request) (*http.Response, error) {
	t.called = true
	return nil, errors.New("transport should not be called")
}

func TestListModelsRejectsUnsafeBaseURLBeforeTransport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		baseURL string
		wantErr string
	}{
		{name: "unsupported scheme", baseURL: "ftp://arena.example", wantErr: "must use http or https"},
		{name: "missing host", baseURL: "https:///arena", wantErr: "must use http or https"},
		{name: "userinfo", baseURL: "https://user:pass@arena.example", wantErr: "must not contain userinfo"},
		{name: "query", baseURL: "https://arena.example?target=other", wantErr: "must not contain query or fragment"},
		{name: "fragment", baseURL: "https://arena.example#fragment", wantErr: "must not contain query or fragment"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			transport := &modelsPreflightTransport{}
			client := NewClient(ClientOptions{
				BaseURL:    tt.baseURL,
				APIKey:     "test-key",
				HTTPClient: &http.Client{Transport: transport},
			})

			_, err := client.ListModels(context.Background())
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("ListModels() error = %v, want error containing %q", err, tt.wantErr)
			}
			if transport.called {
				t.Fatal("ListModels() invoked HTTP transport for rejected base URL")
			}
		})
	}
}
