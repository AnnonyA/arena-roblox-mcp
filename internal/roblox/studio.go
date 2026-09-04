package roblox

import (
	"encoding/json"
	"errors"
)

var (
	ErrNoStudioSessions        = errors.New("no Roblox Studio MCP session detected")
	ErrStudioSelectionRequired = errors.New("multiple Roblox Studio sessions detected; studio_id is required")
	ErrStudioNotFound          = errors.New("requested Roblox Studio session not found")
	ErrStudioIDRequired        = errors.New("Roblox Studio session id is required")
	ErrStudioIDMismatch        = errors.New("tool arguments target a different Roblox Studio session")
)

type StudioSession struct {
	ID      string
	Name    string
	PlaceID string
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
	for _, studio := range response.Studios {
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

	if existing, ok := args["studio_id"]; ok && existing != studioID {
		return nil, ErrStudioIDMismatch
	}
	args["studio_id"] = studioID

	return json.Marshal(args)
}
