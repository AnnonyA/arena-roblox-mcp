package arena

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
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
						return false, fmt.Errorf("conflicting tool call id for index %d", fragment.Index)
					}
				}
				if fragment.Type != "" {
					if call.Type != "" && call.Type != fragment.Type {
						return false, fmt.Errorf("conflicting tool call type for index %d", fragment.Index)
					}
					if fragment.Type != "function" {
						return false, fmt.Errorf("unsupported tool call type %q for index %d", fragment.Type, fragment.Index)
					}
					if call.Type == "" {
						call.Type = fragment.Type
					}
				}
				name := fragment.Function.Name
				switch {
				case name == "":
				case strings.HasPrefix(name, call.Function.Name):
					call.Function.Name = name
				default:
					call.Function.Name += name
				}
				if len(call.Function.Name) > maxStreamToolCallNameBytes {
					return false, fmt.Errorf("tool call name exceeds %d bytes for index %d", maxStreamToolCallNameBytes, fragment.Index)
				}
				beforeArgumentBytes := len(call.Function.Arguments)
				arguments := fragment.Function.Arguments
				switch {
				case arguments == "":
				case strings.HasPrefix(arguments, call.Function.Arguments):
					call.Function.Arguments = arguments
				default:
					call.Function.Arguments += arguments
				}
				if len(call.Function.Arguments) > maxStreamToolCallArgumentBytes {
					return false, fmt.Errorf("tool call arguments exceed %d bytes for index %d", maxStreamToolCallArgumentBytes, fragment.Index)
				}
				addedArgumentBytes := len(call.Function.Arguments) - beforeArgumentBytes
				if addedArgumentBytes > 0 {
					totalToolCallArgumentBytes += addedArgumentBytes
					if totalToolCallArgumentBytes > maxStreamTotalToolCallArgumentBytes {
						return false, fmt.Errorf("total tool call arguments exceed %d bytes", maxStreamTotalToolCallArgumentBytes)
					}
				}
			}
		}
		return false, nil
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), maxSSELineBytes)
	done := false
	firstLine := true
	var dataLines []string
	dataBytes := 0
	for scanner.Scan() {
		line := scanner.Text()
		if !utf8.ValidString(line) {
			return ChatResult{}, fmt.Errorf("Arena stream contains invalid UTF-8")
		}
		if firstLine {
			line = strings.TrimPrefix(line, "\ufeff")
			firstLine = false
		}
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if len(dataLines) >= maxSSEDataLines {
				return ChatResult{}, fmt.Errorf("Arena SSE event has too many data lines: limit %d", maxSSEDataLines)
			}
			added := len(data)
			if len(dataLines) > 0 {
				added++
			}
			if dataBytes+added > maxSSEEventBytes {
				return ChatResult{}, fmt.Errorf("Arena SSE event exceeds %d bytes", maxSSEEventBytes)
			}
			dataLines = append(dataLines, data)
			dataBytes += added
			continue
		}
		if line != "" || len(dataLines) == 0 {
			continue
		}

		payload := strings.Join(dataLines, "\n")
		dataLines = dataLines[:0]
		dataBytes = 0
		done, err = processPayload(payload)
		if err != nil {
			return ChatResult{}, err
		}
		if done {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return ChatResult{}, fmt.Errorf("read Arena stream: %w", err)
	}
	if !done && len(dataLines) > 0 {
		done, err = processPayload(strings.Join(dataLines, "\n"))
		if err != nil {
			return ChatResult{}, err
		}
	}
	if !done {
		return ChatResult{}, fmt.Errorf("Arena stream ended before [DONE]")
	}
	indexes := make([]int, 0, len(calls))
	for index := range calls {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	for _, index := range indexes {
		call := calls[index]
		if strings.TrimSpace(call.ID) == "" {
			return ChatResult{}, fmt.Errorf("missing tool call id for index %d", index)
		}
		if call.Type == "" {
			return ChatResult{}, fmt.Errorf("missing tool call type for index %d", index)
		}
		if strings.TrimSpace(call.Function.Name) == "" {
			return ChatResult{}, fmt.Errorf("missing tool call name for index %d", index)
		}
		if !json.Valid([]byte(call.Function.Arguments)) {
			return ChatResult{}, fmt.Errorf("invalid tool call arguments for index %d", index)
		}
		if !strings.HasPrefix(strings.TrimSpace(call.Function.Arguments), "{") {
			return ChatResult{}, fmt.Errorf("tool call arguments for index %d must be a JSON object", index)
		}
		if err := rejectDuplicateStreamJSONKeys([]byte(call.Function.Arguments)); err != nil {
			return ChatResult{}, fmt.Errorf("invalid tool call arguments for index %d: %w", index, err)
		}
		result.ToolCalls = append(result.ToolCalls, *call)
	}
	return result, nil
}

func retryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || (status >= 500 && status <= 599)
}

func waitRetry(ctx context.Context, retryAfter string) error {
	delay := retryAfterDelay(retryAfter, time.Now())
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func retryAfterDelay(value string, now time.Time) time.Duration {
	value = strings.TrimSpace(value)
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil && seconds >= 0 {
		if seconds > int64((time.Duration(1<<63-1))/time.Second) {
			return defaultRetryDelay
		}
		delay := time.Duration(seconds) * time.Second
		if delay > maxRetryDelay {
			return maxRetryDelay
		}
		return delay
	}
	if at, err := http.ParseTime(value); err == nil {
		if !at.After(now) {
			return 0
		}
		delay := at.Sub(now)
		if delay > maxRetryDelay {
			return maxRetryDelay
		}
		return delay
	}
	return defaultRetryDelay
}