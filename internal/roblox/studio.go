package roblox

import (
	"context"
	"encoding/json"
	"errors"

	arenamcp "github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
)

var (
	ErrNoStudioSessions        = errors.New("no Roblox Studio MCP session detected")
	ErrStudioSelectionRequired = errors.New("multiple Roblox Studio sessions detected; studio_id is required")
	ErrStudioNotFound          = errors.New("requested Roblox Studio session not found")
	ErrStudioIDRequired        = errors.New("Roblox Studio session id is required")
	ErrStudioIDMismatch        = errors.New("tool arguments target a different Roblox Studio session")
	ErrStudioDiscoveryFailed   = errors.New("Roblox Studio MCP discovery failed")
	ErrInvalidStudioSession    = errors.New("Roblox Studio MCP discovery returned a session without studio_id")
	ErrDuplicateStudioSession  = errors.New("Roblox Studio MCP discovery returned duplicate studio_id")
)

type StudioSession struct {
	ID      string
	Name    string
	PlaceID string
}

type ToolCaller interface {
	CallTool(context.Context, string, json.RawMessage) (arenamcp.ToolResult, error)
}

func DiscoverStudioSessions(ctx context.Context, caller ToolCaller) ([]StudioSession, error) {
	result, err := caller.CallTool(ctx, "list_roblox_studios", json.RawMessage(`{}`))
	if err != nil {
		return nil, err
	}
	if result.IsError {
		return nil, ErrStudioDiscoveryFailed
	}
	return ParseStudioSessions(result.StructuredContent)
}

func ParseStudioSessions(payload json.RawMessage) ([]StudioSession, error) {
	var response struct {
		Studios []struct {
			ID      string      `json:"studio_id"`
			Name    string      `json:"name"`
			PlaceID json.Number `json:"place_id"`
		} `json:"studios"`
	}
	if err := json.Unmarshal(payload, &response); err != nil {
		return nil, err
	}

	sessions := make([]StudioSession, 0, len(response.Studios))
	seenIDs := make(map[string]struct{}, len(response.Studios))
	for _, studio := range response.Studios {
		if studio.ID == "" {
			return nil, ErrInvalidStudioSession
		}
		if _, exists := seenIDs[studio.ID]; exists {
			return nil, ErrDuplicateStudioSession
		}
		seenIDs[studio.ID] = struct{}{}
		sessions = append(sessions, StudioSession{
			ID:      studio.ID,
			Name:    studio.Name,
			PlaceID: studio.PlaceID.String(),
		})
	}
	return sessions, nil
}

func SelectStudio(sessions []StudioSession, requestedID string) (StudioSession, error) {
	if len(sessions) == 0 {
		return StudioSession{}, ErrNoStudioSessions
	}

	if requestedID != "" {
		for _, session := range sessions {
			if session.ID == requestedID {
				return session, nil
			}
		}
		return StudioSession{}, ErrStudioNotFound
	}

	if len(sessions) == 1 {
		return sessions[0], nil
	}

	return StudioSession{}, ErrStudioSelectionRequired
}

func TargetStudio(arguments json.RawMessage, studioID string) (json.RawMessage, error) {
	if studioID == "" {
		return nil, ErrStudioIDRequired
	}

	var args map[string]any
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	if args == nil {
		args = make(map[string]any)
	}

	if existing, ok := args["studio_id"]; ok && existing != studioID {
		return nil, ErrStudioIDMismatch
	}
	args["studio_id"] = studioID

	return json.Marshal(args)
}
