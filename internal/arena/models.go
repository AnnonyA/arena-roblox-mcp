package arena

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Model struct {
	ID string `json:"id"`
}

type modelsResponse struct {
	Data []Model `json:"data"`
}

func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	var resp *http.Response
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/models", nil)
		if err != nil {
			return nil, fmt.Errorf("build models request: %w", err)
		}
		req.Header.Set("Accept", "application/json")
		if c.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+c.apiKey)
		}

		resp, err = c.http.Do(req)
		if err != nil {
			return nil, fmt.Errorf("list Arena models: %w", err)
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			break
		}
		if attempt == 2 || !retryableStatus(resp.StatusCode) {
			defer resp.Body.Close()
			return nil, fmt.Errorf("Arena models request failed with HTTP %d", resp.StatusCode)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if err := waitRetry(ctx, resp.Header.Get("Retry-After")); err != nil {
			return nil, fmt.Errorf("wait to retry Arena models: %w", err)
		}
	}
	defer resp.Body.Close()

	var payload modelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode Arena models: %w", err)
	}
	return payload.Data, nil
}
