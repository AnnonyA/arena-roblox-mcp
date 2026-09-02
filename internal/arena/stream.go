package arena

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	Index    int          `json:"index,omitempty"`
	ID       string       `json:"id,omitempty"`
	Type     string       `json:"type,omitempty"`
	Function FunctionCall `json:"function"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages,omitempty"`
	Stream   bool      `json:"stream"`
}

type ChatResult struct {
	Text      string
	ToolCalls []ToolCall
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls"`
		} `json:"delta"`
	} `json:"choices"`
}

func (c *Client) StreamChat(ctx context.Context, req ChatRequest, onText func(string)) (ChatResult, error) {
	req.Stream = true
	body, err := json.Marshal(req)
	if err != nil {
		return ChatResult{}, fmt.Errorf("encode Arena chat request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return ChatResult{}, fmt.Errorf("build Arena chat request: %w", err)
	}
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return ChatResult{}, fmt.Errorf("stream Arena chat: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ChatResult{}, fmt.Errorf("Arena chat request failed with HTTP %d", resp.StatusCode)
	}

	var result ChatResult
	calls := map[int]*ToolCall{}
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" {
			continue
		}
		if payload == "[DONE]" {
			break
		}

		var chunk streamChunk
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			return ChatResult{}, fmt.Errorf("decode Arena stream chunk: %w", err)
		}
		for _, choice := range chunk.Choices {
			if choice.Delta.Content != "" {
				result.Text += choice.Delta.Content
				if onText != nil {
					onText(choice.Delta.Content)
				}
			}
			for _, fragment := range choice.Delta.ToolCalls {
				call := calls[fragment.Index]
				if call == nil {
					call = &ToolCall{Index: fragment.Index}
					calls[fragment.Index] = call
				}
				call.ID += fragment.ID
				call.Type += fragment.Type
				call.Function.Name += fragment.Function.Name
				call.Function.Arguments += fragment.Function.Arguments
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ChatResult{}, fmt.Errorf("read Arena stream: %w", err)
	}
	for i := 0; i < len(calls); i++ {
		if call := calls[i]; call != nil {
			result.ToolCalls = append(result.ToolCalls, *call)
		}
	}
	return result, nil
}
