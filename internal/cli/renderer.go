package cli

import "fmt"

const maxStatusDisplayRunes = 256

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
	arena := boundedSafeDisplayText(status.Arena, maxStatusDisplayRunes)
	mcpStatus := boundedSafeDisplayText(status.MCP, maxStatusDisplayRunes)
	studio := boundedSafeDisplayText(status.Studio, maxStatusDisplayRunes)
	model := boundedSafeDisplayText(status.Model, maxStatusDisplayRunes)
	session := boundedSafeDisplayText(status.Session, maxStatusDisplayRunes)

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
