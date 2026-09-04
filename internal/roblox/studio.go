package roblox

import "errors"

var (
	ErrNoStudioSessions        = errors.New("no Roblox Studio MCP session detected")
	ErrStudioSelectionRequired = errors.New("multiple Roblox Studio sessions detected; studio_id is required")
	ErrStudioNotFound          = errors.New("requested Roblox Studio session not found")
)

type StudioSession struct {
	ID   string
	Name string
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
