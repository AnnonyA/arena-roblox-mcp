package cli

import "fmt"

type StartupStatus struct {
	Arena   string
	Studio  string
	Model   string
	Session string
}

func StartupText(status StartupStatus) string {
	return fmt.Sprintf(
		"Arena Roblox MCP\n────────────────────────────\nArena      %s\nStudio     %s\nModel      %s\nSession    %s\n\n> ",
		status.Arena,
		status.Studio,
		status.Model,
		status.Session,
	)
}
