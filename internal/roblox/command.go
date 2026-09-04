package roblox

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
)

func DefaultWindowsCommand(localAppData string) (mcp.CommandConfig, error) {
	if strings.TrimSpace(localAppData) == "" {
		return mcp.CommandConfig{}, errors.New("LOCALAPPDATA is required for Roblox Studio MCP")
	}

	return mcp.CommandConfig{
		Command: "cmd.exe",
		Args: []string{
			"/c",
			filepath.Join(localAppData, "Roblox", "mcp.bat"),
		},
	}, nil
}
