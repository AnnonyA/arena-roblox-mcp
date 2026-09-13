package agent

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

const (
	maxToolCallsPerRound  = 32
	maxToolCallIDBytes    = 256
	maxToolCallNameBytes  = 256
	maxToolArgumentDepth  = 64
	maxToolArgumentBytes  = 1 << 20
	maxToolArgumentValues = 16384
	maxToolResultBytes    = 256 << 10

	toolResultTruncationMarker = "… [tool result truncated]"
)

var (
	ErrNoChatRound      = errors.New("chat round is not configured")
	ErrNoToolDispatcher = errors.New("tool dispatcher is not configured")
	ErrInvalidToolCall  = errors.New("invalid tool call")
)

type ChatRoundFunc func(context.Context, []arena.Message) (arena.ChatResult, error)

func RunConversation(ctx context.Context, maxRounds int, messages []arena.Message, runRound ChatRoundFunc, dispatcher *ToolDispatcher) (string, error) {
	if runRound == nil {
		return "", ErrNoChatRound
	}
	if maxRounds <= 0 {
		maxRounds = DefaultMaxToolRounds
	}

	conversation := append([]arena.Message(nil), messages...)
	var finalText string
	toolRounds := 0
	err := RunToolLoop(ctx, maxRounds+1, func(ctx context.Context) (bool, error) {
		result, err := runRound(ctx, conversation)
		if err != nil {
			return false, err
		}
		if len(result.ToolCalls) == 0 {
			finalText = result.Text
			return false, nil
		}
		if len(result.ToolCalls) > maxToolCallsPerRound {
			return false, ErrInvalidToolCall
		}
		if toolRounds >= maxRounds {
			return false, ErrMaxToolRounds
		}
		if dispatcher == nil {
			return false, ErrNoToolDispatcher
		}
		seenCallIDs := make(map[string]struct{}, len(result.ToolCalls))
		for _, call := range result.ToolCalls {
			trimmedID := strings.TrimSpace(call.ID)
			trimmedName := strings.TrimSpace(call.Function.Name)
			if call.Type != "function" || trimmedID == "" || trimmedID != call.ID || len(call.ID) > maxToolCallIDBytes || strings.IndexFunc(call.ID, unicode.IsControl) >= 0 || strings.IndexFunc(call.ID, func(r rune) bool { return unicode.Is(unicode.Cf, r) }) >= 0 || trimmedName == "" || trimmedName != call.Function.Name || len(call.Function.Name) > maxToolCallNameBytes || strings.IndexFunc(call.Function.Name, unicode.IsControl) >= 0 || strings.IndexFunc(call.Function.Name, func(r rune) bool { return unicode.Is(unicode.Cf, r) }) >= 0 || !validToolArguments(call.Function.Arguments) {
				return false, ErrInvalidToolCall
			}
			if _, exists := seenCallIDs[call.ID]; exists {
				return false, ErrInvalidToolCall
			}
			seenCallIDs[call.ID] = struct{}{}
		}
		toolRounds++

		conversation = append(conversation, arena.Message{
			Role:      "assistant",
			Content:   result.Text,
			ToolCalls: append([]arena.ToolCall(nil), result.ToolCalls...),
		})
		for _, call := range result.ToolCalls {
			toolResult, err := dispatcher.Dispatch(ctx, call.Function.Name, json.RawMessage(call.Function.Arguments))
			var content json.RawMessage
			if err != nil {
				if ctxErr := ctx.Err(); ctxErr != nil {
					return false, ctxErr
				}
				content, err = json.Marshal(map[string]string{"error": err.Error()})
				if err != nil {
					return false, err
				}
			} else {
				content = toolResult.StructuredContent
				if len(content) == 0 {
					content = toolResult.Content
				}
				if len(content) == 0 {
					content = json.RawMessage("null")
				}
			}
			conversation = append(conversation, arena.Message{
				Role:       "tool",
				Content:    compactToolResult(content),
				ToolCallID: call.ID,
			})
		}
		return true, nil
	})
	if err != nil {
		return "", err
	}
	return finalText, nil
}

func compactToolResult(content json.RawMessage) string {
	contentText := string(content)
	if !utf8.ValidString(contentText) {
		contentText = strings.ToValidUTF8(contentText, "\uFFFD")
	}
	if len(contentText) <= maxToolResultBytes {
		return contentText
	}

	limit := maxToolResultBytes - len(toolResultTruncationMarker)
	for limit > 0 && !utf8.RuneStart(contentText[limit]) {
		limit--
	}
	return contentText[:limit] + toolResultTruncationMarker
}

func validToolArguments(arguments string) bool {
	if len(arguments) > maxToolArgumentBytes || !utf8.ValidString(arguments) {
		return false
	}
	decoder := json.NewDecoder(strings.NewReader(arguments))
	values := 0
	if !validJSONObject(decoder, &values) {
		return false
	}
	_, err := decoder.Token()
	return errors.Is(err, io.EOF)
}

func validJSONObject(decoder *json.Decoder, values *int) bool {
	start, err := decoder.Token()
	if err != nil || start != json.Delim('{') {
		return false
	}
	return validJSONObjectContents(decoder, 1, values)
}

func validJSONObjectContents(decoder *json.Decoder, depth int, values *int) bool {
	seen := make(map[string]struct{})
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return false
		}
		key, ok := token.(string)
		if !ok {
			return false
		}
		if _, exists := seen[key]; exists {
			return false
		}
		seen[key] = struct{}{}
		if !consumeToolArgumentValue(values) {
			return false
		}

		if !validJSONValue(decoder, depth, values) {
			return false
		}
	}
	end, err := decoder.Token()
	return err == nil && end == json.Delim('}')
}

func validJSONValue(decoder *json.Decoder, depth int, values *int) bool {
	token, err := decoder.Token()
	if err != nil {
		return false
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return true
	}

	nextDepth := depth + 1
	if nextDepth > maxToolArgumentDepth {
		return false
	}

	switch delim {
	case '{':
		return validJSONObjectContents(decoder, nextDepth, values)
	case '[':
		for decoder.More() {
			if !consumeToolArgumentValue(values) || !validJSONValue(decoder, nextDepth, values) {
				return false
			}
		}
		end, err := decoder.Token()
		return err == nil && end == json.Delim(']')
	default:
		return false
	}
}

func consumeToolArgumentValue(values *int) bool {
	*values++
	return *values <= maxToolArgumentValues
}
