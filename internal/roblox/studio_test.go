package roblox

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestSelectStudioRejectsNoSessions(t *testing.T) {
	_, err := SelectStudio(nil, "")
	if !errors.Is(err, ErrNoStudioSessions) {
		t.Fatalf("SelectStudio() error = %v, want ErrNoStudioSessions", err)
	}
}

func TestSelectStudioAutomaticallyUsesOnlySession(t *testing.T) {
	sessions := []StudioSession{{ID: "studio-1", Name: "Main"}}

	got, err := SelectStudio(sessions, "")
	if err != nil {
		t.Fatalf("SelectStudio() error = %v", err)
	}
	if got != sessions[0] {
		t.Fatalf("SelectStudio() = %#v, want %#v", got, sessions[0])
	}
}

func TestSelectStudioRequiresExplicitIDWhenMultipleSessionsExist(t *testing.T) {
	sessions := []StudioSession{
		{ID: "studio-1", Name: "Main"},
		{ID: "studio-2", Name: "Test"},
	}

	_, err := SelectStudio(sessions, "")
	if !errors.Is(err, ErrStudioSelectionRequired) {
		t.Fatalf("SelectStudio() error = %v, want ErrStudioSelectionRequired", err)
	}
}

func TestSelectStudioUsesRequestedID(t *testing.T) {
	sessions := []StudioSession{
		{ID: "studio-1", Name: "Main"},
		{ID: "studio-2", Name: "Test"},
	}

	got, err := SelectStudio(sessions, "studio-2")
	if err != nil {
		t.Fatalf("SelectStudio() error = %v", err)
	}
	if got != sessions[1] {
		t.Fatalf("SelectStudio() = %#v, want %#v", got, sessions[1])
	}
}

func TestSelectStudioRejectsUnknownRequestedID(t *testing.T) {
	sessions := []StudioSession{{ID: "studio-1", Name: "Main"}}

	_, err := SelectStudio(sessions, "missing")
	if !errors.Is(err, ErrStudioNotFound) {
		t.Fatalf("SelectStudio() error = %v, want ErrStudioNotFound", err)
	}
}

func TestTargetStudioAddsSelectedStudioID(t *testing.T) {
	got, err := TargetStudio(json.RawMessage(`{"path":"Workspace"}`), "studio-2")
	if err != nil {
		t.Fatalf("TargetStudio() error = %v", err)
	}

	var args map[string]any
	if err := json.Unmarshal(got, &args); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if args["studio_id"] != "studio-2" {
		t.Fatalf("studio_id = %#v, want studio-2", args["studio_id"])
	}
	if args["path"] != "Workspace" {
		t.Fatalf("path = %#v, want Workspace", args["path"])
	}
}

func TestTargetStudioRejectsMissingStudioID(t *testing.T) {
	_, err := TargetStudio(json.RawMessage(`{"path":"Workspace"}`), "")
	if !errors.Is(err, ErrStudioIDRequired) {
		t.Fatalf("TargetStudio() error = %v, want ErrStudioIDRequired", err)
	}
}

func TestTargetStudioRejectsConflictingStudioID(t *testing.T) {
	_, err := TargetStudio(json.RawMessage(`{"studio_id":"studio-1"}`), "studio-2")
	if !errors.Is(err, ErrStudioIDMismatch) {
		t.Fatalf("TargetStudio() error = %v, want ErrStudioIDMismatch", err)
	}
}
