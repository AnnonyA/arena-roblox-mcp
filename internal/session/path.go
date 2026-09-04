package session

import (
	"os"
	"path/filepath"
)

// DataDir returns the directory used for persisted arena-rbx session data.
func DataDir() (string, error) {
	return resolveDataDir(os.Getenv("LOCALAPPDATA"), os.UserConfigDir)
}

func resolveDataDir(localAppData string, userConfigDir func() (string, error)) (string, error) {
	base := localAppData
	if base == "" {
		fallback, err := userConfigDir()
		if err != nil {
			return "", err
		}
		base = fallback
	}

	return filepath.Join(base, "arena-rbx", "sessions"), nil
}
