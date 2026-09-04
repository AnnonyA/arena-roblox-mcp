package roblox

import (
	"path/filepath"

	"github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
)

func DefaultWindowsCommand(localAppData string) (mcp.CommandConfig, error) {
	return mcp.CommandConfig{
		Command: "cmd.exe",
		Args: []string{
			"/c",
			filepath.Join(localAppData, "Roblox", "mcp.bat"),
		},
	}, nil
}
