package arena

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestListModelsRetriesTransportFailure(t *testing.T) {
	t.Parallel()

	attempts := 0
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return nil, errors.New("temporary network failure")
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"data":[{"id":"model-a"}]}`)),
			Request:    req,
		}, nil
	})}

	c := NewClient(ClientOptions{BaseURL: "https://arena.invalid", HTTPClient: client})
	got, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
	if len(got) != 1 || got[0].ID != "model-a" {
		t.Fatalf("models = %#v", got)
	}
}
