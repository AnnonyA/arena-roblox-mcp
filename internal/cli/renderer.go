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
	mcp := ""
	if status.MCP != "" {
		mcp = fmt.Sprintf("MCP        %s\n", status.MCP)
	}
	return fmt.Sprintf(
		"Arena      %s\n%sStudio     %s\nModel      %s\nSession    %s\n",
		status.Arena,
		mcp,
		status.Studio,
		status.Model,
		status.Session,
	)
}
