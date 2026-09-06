package cli

import "fmt"

type StartupStatus struct {
	Arena   string
	MCP     string
	Studio  string
	Model   string
	Session string
}

func StartupText(status StartupStatus) string {
	return fmt.Sprintf(
		"Arena Roblox MCP\n────────────────────────────\n%s\n> ",
		StatusText(status),
	)
}

func StatusText(status StartupStatus) string {
	return fmt.Sprintf(
		"Arena      %s\nMCP        %s\nStudio     %s\nModel      %s\nSession    %s\n",
		status.Arena,
		status.MCP,
		status.Studio,
		status.Model,
		status.Session,
	)
}
