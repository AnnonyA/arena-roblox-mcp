package arena

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxModelsResponseBytes = 1024 * 1024

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
			if attempt == 2 || ctx.Err() != nil {
				return nil, fmt.Errorf("list Arena models: %w", err)
			}
			if err := waitRetry(ctx, ""); err != nil {
				return nil, fmt.Errorf("wait to retry Arena models: %w", err)
			}
			continue
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			break
		}
		if attempt == 2 || !retryableStatus(resp.StatusCode) {
			drainAndClose(resp.Body)
			return nil, arenaStatusError("models", resp.StatusCode)
		}
		drainAndClose(resp.Body)
		if err := waitRetry(ctx, resp.Header.Get("Retry-After")); err != nil {
			return nil, fmt.Errorf("wait to retry Arena models: %w", err)
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxModelsResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read Arena models: %w", err)
	}
	if len(body) > maxModelsResponseBytes {
		return nil, fmt.Errorf("Arena models response exceeds %d bytes", maxModelsResponseBytes)
	}

	var payload modelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode Arena models: %w", err)
	}
	seen := make(map[string]struct{}, len(payload.Data))
	for i, model := range payload.Data {
		if strings.TrimSpace(model.ID) == "" {
			return nil, fmt.Errorf("decode Arena models: blank model id at index %d", i)
		}
		if _, ok := seen[model.ID]; ok {
			return nil, fmt.Errorf("decode Arena models: duplicate model id %q at index %d", model.ID, i)
		}
		seen[model.ID] = struct{}{}
	}
	return payload.Data, nil
}
