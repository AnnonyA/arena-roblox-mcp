package arena

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	defaultRetryDelay                   = 200 * time.Millisecond
	maxRetryDelay                       = time.Minute
	maxSSEEventBytes                    = 1024 * 1024
	maxSSELineBytes                     = maxSSEEventBytes + len("data: ") + 1
	maxSSEDataLines                     = 4096
	maxStreamTextBytes                  = maxSSEEventBytes
	maxStreamToolCalls                  = 128
	maxStreamToolCallIDBytes            = 4096
	maxStreamToolCallNameBytes          = 4096
	maxStreamToolCallArgumentBytes      = maxSSEEventBytes
	maxStreamTotalToolCallArgumentBytes = maxSSEEventBytes
)

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
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
	Model    string           `json:"model"`
	Messages []Message        `json:"messages,omitempty"`
	Tools    []ToolDefinition `json:"tools,omitempty"`
	Stream   bool             `json:"stream"`
}
type ChatResult struct {
	Text      string
	ToolCalls []ToolCall
}
type streamChunk struct {
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls"`
		} `json:"delta"`
	} `json:"choices"`
}

func (c *Client) StreamChat(ctx context.Context, req ChatRequest, onText func(string)) (ChatResult, error) {
	if c == nil {
		return ChatResult{}, fmt.Errorf("stream Arena chat: client is nil")
	}
	if ctx == nil {
		return ChatResult{}, fmt.Errorf("stream Arena chat: context is nil")
	}
	if c.http == nil {
		return ChatResult{}, fmt.Errorf("stream Arena chat: HTTP client is nil")
	}
	if strings.TrimSpace(c.baseURL) == "" {
		return ChatResult{}, fmt.Errorf("stream Arena chat: base URL is empty")
	}
	baseURL, err := url.Parse(c.baseURL)
	if err != nil || (baseURL.Scheme != "http" && baseURL.Scheme != "https") || baseURL.Host == "" {
		return ChatResult{}, fmt.Errorf("stream Arena chat: base URL must use http or https")
	}
	if baseURL.User != nil {
		return ChatResult{}, fmt.Errorf("stream Arena chat: base URL must not contain userinfo")
	}
	if baseURL.RawQuery != "" || baseURL.Fragment != "" {
		return ChatResult{}, fmt.Errorf("stream Arena chat: base URL must not contain query or fragment")
	}

	req.Stream = true
	body, err := json.Marshal(req)
	if err != nil {
		return ChatResult{}, fmt.Errorf("encode Arena chat request: %w", err)
	}

	var resp *http.Response
	for attempt := 0; attempt < 3; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
		if err != nil {
			return ChatResult{}, fmt.Errorf("build Arena chat request: %w", err)
		}
		httpReq.Header.Set("Accept", "text/event-stream")
		httpReq.Header.Set("Content-Type", "application/json")
		if c.apiKey != "" {
			httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		}

		resp, err = c.http.Do(httpReq)
		if err != nil {
			if attempt == 2 || ctx.Err() != nil {
				return ChatResult{}, fmt.Errorf("stream Arena chat: %w", err)
			}
			if err := waitRetry(ctx, ""); err != nil {
				return ChatResult{}, fmt.Errorf("wait to retry Arena chat: %w", err)
			}
			continue
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			break
		}
		if attempt == 2 || !retryableStatus(resp.StatusCode) {
			drainAndClose(resp.Body)
			return ChatResult{}, arenaStatusError("chat", resp.StatusCode)
		}
		drainAndClose(resp.Body)
		if err := waitRetry(ctx, resp.Header.Get("Retry-After")); err != nil {
			return ChatResult{}, fmt.Errorf("wait to retry Arena chat: %w", err)
		}
	}
	defer resp.Body.Close()
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil || mediaType != "text/event-stream" {
			return ChatResult{}, fmt.Errorf("unexpected Arena chat content type %q: want text/event-stream", contentType)
		}
	}

	var result ChatResult
	calls := map[int]*ToolCall{}
	callIDs := map[string]int{}
	totalToolCallArgumentBytes := 0
	processPayload := func(payload string) (bool, error) {
		if payload == "[DONE]" {
			return true, nil
		}
		if payload == "" {
			return false, nil
		}

		var chunk streamChunk
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			return false, fmt.Errorf("decode Arena stream chunk: %w", err)
		}
		for _, choice := range chunk.Choices {
			if choice.Index < 0 {
				return false, fmt.Errorf("invalid negative choice index %d", choice.Index)
			}
			if choice.Index != 0 {
				continue
			}
			if choice.Delta.Content != "" {
				if len(result.Text)+len(choice.Delta.Content) > maxStreamTextBytes {
					return false, fmt.Errorf("streamed text exceeds %d bytes", maxStreamTextBytes)
				}
				result.Text += choice.Delta.Content
				if onText != nil {
					onText(choice.Delta.Content)
				}
			}
			for _, fragment := range choice.Delta.ToolCalls {
				if fragment.Index < 0 {
					return false, fmt.Errorf("invalid negative tool call index %d", fragment.Index)
				}
				call := calls[fragment.Index]
				if call == nil {
					if len(calls) >= maxStreamToolCalls {
						return false, fmt.Errorf("too many tool calls in Arena stream: limit %d", maxStreamToolCalls)
					}
					call = &ToolCall{Index: fragment.Index}
					calls[fragment.Index] = call
				}
				if fragment.ID != "" {
					if len(fragment.ID) > maxStreamToolCallIDBytes {
						return false, fmt.Errorf("tool call id exceeds %d bytes for index %d", maxStreamToolCallIDBytes, fragment.Index)
					}
					if index, ok := callIDs[fragment.ID]; ok && index != fragment.Index {
						return false, fmt.Errorf("duplicate tool call id %q for indexes %d and %d", fragment.ID, index, fragment.Index)
					}
					if call.ID == "" {
						call.ID = fragment.ID
						callIDs[fragment.ID] = fragment.Index
					} else if call.ID != fragment.ID {
						return false, fmt.Errorf("conflicting tool call id for index %d: %q then %q", fragment.Index, call.ID, fragment.ID)
					}
				}
				if fragment.Type != "" {
					if fragment.Type != "function" {
						return false, fmt.Errorf("unsupported tool call type %q for index %d", fragment.Type, fragment.Index)
					}
					if call.Type == "" {
						call.Type = fragment.Type
					} else if call.Type != fragment.Type {
						return false, fmt.Errorf("conflicting tool call type for index %d: %q then %q", fragment.Index, call.Type, fragment.Type)
					}
				}
				if fragment.Function.Name != "" {
					if len(fragment.Function.Name) > maxStreamToolCallNameBytes {
						return false, fmt.Errorf("tool call name exceeds %d bytes for index %d", maxStreamToolCallNameBytes, fragment.Index)
					}
					if err := validateStreamToolName(fragment.Function.Name); err != nil {
						return false, fmt.Errorf("invalid tool call name for index %d: %w", fragment.Index, err)
					}
					if call.Function.Name == "" {
						call.Function.Name = fragment.Function.Name
					} else if call.Function.Name != fragment.Function.Name {
						if strings.HasPrefix(fragment.Function.Name, call.Function.Name) {
							call.Function.Name = fragment.Function.Name
						} else if !strings.HasPrefix(call.Function.Name, fragment.Function.Name) {
							return false, fmt.Errorf("conflicting tool call name for index %d: %q then %q", fragment.Index, call.Function.Name, fragment.Function.Name)
						}
					}
				}
				if fragment.Function.Arguments != "" {
					if len(fragment.Function.Arguments) > maxStreamToolCallArgumentBytes {
						return false, fmt.Errorf("tool call arguments exceed %d bytes for index %d", maxStreamToolCallArgumentBytes, fragment.Index)
					}
					arguments := fragment.Function.Arguments
					if call.Function.Arguments != "" && strings.HasPrefix(arguments, call.Function.Arguments) {
						arguments = strings.TrimPrefix(arguments, call.Function.Arguments)
					}
					if len(call.Function.Arguments)+len(arguments) > maxStreamToolCallArgumentBytes {
						return false, fmt.Errorf("tool call arguments exceed %d bytes for index %d", maxStreamToolCallArgumentBytes, fragment.Index)
					}
					if totalToolCallArgumentBytes+len(arguments) > maxStreamTotalToolCallArgumentBytes {
						return false, fmt.Errorf("total tool call arguments exceed %d bytes", maxStreamTotalToolCallArgumentBytes)
					}
					call.Function.Arguments += arguments
					totalToolCallArgumentBytes += len(arguments)
				}
			}
		}
		return false, nil
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 4096), maxSSELineBytes)
	var dataLines []string
	eventBytes := 0
	flushEvent := func() (bool, error) {
		if len(dataLines) == 0 {
			return false, nil
		}
		payload := strings.Join(dataLines, "\n")
		dataLines = dataLines[:0]
		eventBytes = 0
		return processPayload(payload)
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			done, err := flushEvent()
			if err != nil {
				return ChatResult{}, err
			}
			if done {
				break
			}
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimPrefix(line, "data:")
		if strings.HasPrefix(data, " ") {
			data = data[1:]
		}
		if len(dataLines) >= maxSSEDataLines {
			return ChatResult{}, fmt.Errorf("Arena stream event exceeds %d data lines", maxSSEDataLines)
		}
		if eventBytes+len(data) > maxSSEEventBytes {
			return ChatResult{}, fmt.Errorf("Arena stream event exceeds %d bytes", maxSSEEventBytes)
		}
		dataLines = append(dataLines, data)
		eventBytes += len(data)
	}
	if err := scanner.Err(); err != nil {
		return ChatResult{}, fmt.Errorf("read Arena stream: %w", err)
	}
	if len(dataLines) > 0 {
		if _, err := flushEvent(); err != nil {
			return ChatResult{}, err
		}
	}

	indexes := make([]int, 0, len(calls))
	for index := range calls {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	for _, index := range indexes {
		call := calls[index]
		if call.ID == "" {
			return ChatResult{}, fmt.Errorf("tool call missing id for index %d", index)
		}
		if call.Type == "" {
			return ChatResult{}, fmt.Errorf("tool call missing type for index %d", index)
		}
		if call.Function.Name == "" {
			return ChatResult{}, fmt.Errorf("tool call missing function name for index %d", index)
		}
		if err := validateToolCallArguments(call.Function.Arguments); err != nil {
			return ChatResult{}, fmt.Errorf("invalid tool call arguments for index %d: %w", index, err)
		}
		result.ToolCalls = append(result.ToolCalls, *call)
	}
	return result, nil
}

func retryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= 500
}

func waitRetry(ctx context.Context, retryAfter string) error {
	delay := defaultRetryDelay
	if retryAfter != "" {
		if seconds, err := strconv.Atoi(strings.TrimSpace(retryAfter)); err == nil && seconds >= 0 {
			delay = time.Duration(seconds) * time.Second
			if delay > maxRetryDelay {
				delay = maxRetryDelay
			}
		}
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func validateStreamToolName(name string) error {
	if !utf8.ValidString(name) {
		return fmt.Errorf("invalid UTF-8")
	}
	if strings.TrimSpace(name) != name {
		return fmt.Errorf("surrounding whitespace")
	}
	for _, r := range name {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("control character")
		}
		if r == '\u2028' || r == '\u2029' {
			return fmt.Errorf("line separator")
		}
		if r >= '\u202a' && r <= '\u202e' || r >= '\u2066' && r <= '\u2069' {
			return fmt.Errorf("bidirectional control character")
		}
	}
	return nil
}
