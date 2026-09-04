package roblox

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefaultWindowsCommandUsesLocalAppDataRobloxMCPBatch(t *testing.T) {
	localAppData := filepath.Join(`C:\Users\tester`, "AppData", "Local")

	got, err := DefaultWindowsCommand(localAppData)
	if err != nil {
		t.Fatalf("DefaultWindowsCommand() error = %v", err)
	}

	if got.Command != "cmd.exe" {
		t.Fatalf("Command = %q, want cmd.exe", got.Command)
	}
	wantArgs := []string{"/c", filepath.Join(localAppData, "Roblox", "mcp.bat")}
	if !reflect.DeepEqual(got.Args, wantArgs) {
		t.Fatalf("Args = %#v, want %#v", got.Args, wantArgs)
	}
}
