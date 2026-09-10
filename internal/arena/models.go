package arena

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxModelsResponseBytes = 1024 * 1024
	maxModelIDBytes        = 4096
	maxModelCount          = 4096
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
	if !utf8.Valid(body) {
		return nil, fmt.Errorf("decode Arena models: invalid UTF-8")
	}

	var payload modelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode Arena models: %w", err)
	}
	if len(payload.Data) > maxModelCount {
		return nil, fmt.Errorf("decode Arena models: model count exceeds %d", maxModelCount)
	}
	seen := make(map[string]struct{}, len(payload.Data))
	for i, model := range payload.Data {
		trimmedID := strings.TrimSpace(model.ID)
		if trimmedID == "" {
			return nil, fmt.Errorf("decode Arena models: blank model id at index %d", i)
		}
		if trimmedID != model.ID {
			return nil, fmt.Errorf("decode Arena models: model id has surrounding whitespace at index %d", i)
		}
		if len(model.ID) > maxModelIDBytes {
			return nil, fmt.Errorf("decode Arena models: model id exceeds %d bytes at index %d", maxModelIDBytes, i)
		}
		for _, r := range model.ID {
			if unicode.IsControl(r) {
				return nil, fmt.Errorf("decode Arena models: model id contains control character at index %d", i)
			}
			if isBidirectionalFormatting(r) {
				return nil, fmt.Errorf("decode Arena models: model id contains bidirectional formatting at index %d", i)
			}
		}
		if _, ok := seen[model.ID]; ok {
			return nil, fmt.Errorf("decode Arena models: duplicate model id %q at index %d", model.ID, i)
		}
		seen[model.ID] = struct{}{}
	}
	return payload.Data, nil
}

func isBidirectionalFormatting(r rune) bool {
	return r >= '\u202a' && r <= '\u202e' || r >= '\u2066' && r <= '\u2069'
}
