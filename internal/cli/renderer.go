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
	arena := safeDisplayText(status.Arena)
	mcpStatus := safeDisplayText(status.MCP)
	studio := safeDisplayText(status.Studio)
	model := safeDisplayText(status.Model)
	session := safeDisplayText(status.Session)

	mcp := ""
	if mcpStatus != "" {
		mcp = fmt.Sprintf("MCP        %s\n", mcpStatus)
	}
	return fmt.Sprintf(
		"Arena      %s\n%sStudio     %s\nModel      %s\nSession    %s\n",
		arena,
		mcp,
		studio,
		model,
		session,
	)
}
