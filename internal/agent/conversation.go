package agent

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/AnnonyA/arena-roblox-mcp/internal/arena"
)

var (
	ErrNoChatRound      = errors.New("chat round is not configured")
	ErrNoToolDispatcher = errors.New("tool dispatcher is not configured")
)

type ChatRoundFunc func(context.Context, []arena.Message) (arena.ChatResult, error)

func RunConversation(ctx context.Context, maxRounds int, messages []arena.Message, runRound ChatRoundFunc, dispatcher *ToolDispatcher) (string, error) {
	if runRound == nil {
		return "", ErrNoChatRound
	}

	conversation := append([]arena.Message(nil), messages...)
	var finalText string
	err := RunToolLoop(ctx, maxRounds, func(ctx context.Context) (bool, error) {
		result, err := runRound(ctx, conversation)
		if err != nil {
			return false, err
		}
		if len(result.ToolCalls) == 0 {
			finalText = result.Text
			return false, nil
		}
		if dispatcher == nil {
			return false, ErrNoToolDispatcher
		}

		conversation = append(conversation, arena.Message{
			Role:      "assistant",
			Content:   result.Text,
			ToolCalls: append([]arena.ToolCall(nil), result.ToolCalls...),
		})
		for _, call := range result.ToolCalls {
			toolResult, err := dispatcher.Dispatch(ctx, call.Function.Name, json.RawMessage(call.Function.Arguments))
			if err != nil {
				return false, err
			}
			content := toolResult.StructuredContent
			if len(content) == 0 {
				content = toolResult.Content
			}
			if len(content) == 0 {
				content = json.RawMessage("null")
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
