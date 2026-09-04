package roblox

import (
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
