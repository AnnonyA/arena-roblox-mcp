package roblox

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestTargetStudioRejectsPaddedStudioID(t *testing.T) {
	_, err := TargetStudio(json.RawMessage(`{"path":"Workspace"}`), " studio-2 ")
	if !errors.Is(err, ErrStudioIDRequired) {
		t.Fatalf("TargetStudio() error = %v, want ErrStudioIDRequired", err)
	}
}
