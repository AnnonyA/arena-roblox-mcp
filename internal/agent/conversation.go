package agent

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"unicode"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
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
			if trimmedID == "" || trimmedID != call.ID || strings.IndexFunc(call.ID, unicode.IsControl) >= 0 || strings.IndexFunc(call.ID, func(r rune) bool { return unicode.Is(unicode.Cf, r) }) >= 0 || trimmedName == "" || trimmedName != call.Function.Name || strings.IndexFunc(call.Function.Name, unicode.IsControl) >= 0 || strings.IndexFunc(call.Function.Name, func(r rune) bool { return unicode.Is(unicode.Cf, r) }) >= 0 || !validToolArguments(call.Function.Arguments) {
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
				Content:    string(content),
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

func validToolArguments(arguments string) bool {
	var object map[string]any
	if err := json.Unmarshal([]byte(arguments), &object); err != nil {
		return false
	}
	return object != nil
}
