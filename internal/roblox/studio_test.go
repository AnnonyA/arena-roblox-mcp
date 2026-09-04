package roblox

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	arenamcp "github.com/AnnonyA/arena-roblox-mcp/internal/mcp"
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

func TestTargetStudioTreatsNullArgumentsAsEmptyObject(t *testing.T) {
	got, err := TargetStudio(json.RawMessage(`null`), "studio-2")
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
}

func TestTargetStudioRejectsMissingStudioID(t *testing.T) {
	_, err := TargetStudio(json.RawMessage(`{"path":"Workspace"}`), "")
	if !errors.Is(err, ErrStudioIDRequired) {
		t.Fatalf("TargetStudio() error = %v, want ErrStudioIDRequired", err)
	}
}

func TestTargetStudioRejectsWhitespaceStudioID(t *testing.T) {
	_, err := TargetStudio(json.RawMessage(`{"path":"Workspace"}`), "   ")
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

func TestParseStudioSessionsReadsListRobloxStudiosPayload(t *testing.T) {
	payload := json.RawMessage(`{"studios":[{"studio_id":"studio-1","name":"Dungeon","place_id":987654},{"studio_id":"studio-local","name":"Local File"}]}`)

	got, err := ParseStudioSessions(payload)
	if err != nil {
		t.Fatalf("ParseStudioSessions() error = %v", err)
	}

	want := []StudioSession{
		{ID: "studio-1", Name: "Dungeon", PlaceID: "987654"},
		{ID: "studio-local", Name: "Local File"},
	}
	if len(got) != len(want) {
		t.Fatalf("len(ParseStudioSessions()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ParseStudioSessions()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestParseStudioSessionsRejectsMissingStudioID(t *testing.T) {
	payload := json.RawMessage(`{"studios":[{"name":"Broken"}]}`)

	_, err := ParseStudioSessions(payload)
	if !errors.Is(err, ErrInvalidStudioSession) {
		t.Fatalf("ParseStudioSessions() error = %v, want ErrInvalidStudioSession", err)
	}
}

func TestParseStudioSessionsRejectsWhitespaceStudioID(t *testing.T) {
	payload := json.RawMessage(`{"studios":[{"studio_id":"   ","name":"Broken"}]}`)

	_, err := ParseStudioSessions(payload)
	if !errors.Is(err, ErrInvalidStudioSession) {
		t.Fatalf("ParseStudioSessions() error = %v, want ErrInvalidStudioSession", err)
	}
}

func TestParseStudioSessionsRejectsDuplicateStudioID(t *testing.T) {
	payload := json.RawMessage(`{"studios":[{"studio_id":"studio-1","name":"First"},{"studio_id":"studio-1","name":"Second"}]}`)

	_, err := ParseStudioSessions(payload)
	if !errors.Is(err, ErrDuplicateStudioSession) {
		t.Fatalf("ParseStudioSessions() error = %v, want ErrDuplicateStudioSession", err)
	}
}

type fakeToolCaller struct {
	name      string
	arguments json.RawMessage
	result    arenamcp.ToolResult
	err       error
}

func (f *fakeToolCaller) CallTool(_ context.Context, name string, arguments json.RawMessage) (arenamcp.ToolResult, error) {
	f.name = name
	f.arguments = append(json.RawMessage(nil), arguments...)
	return f.result, f.err
}

func TestDiscoverStudioSessionsCallsListRobloxStudios(t *testing.T) {
	caller := &fakeToolCaller{result: arenamcp.ToolResult{
		StructuredContent: json.RawMessage(`{"studios":[{"studio_id":"studio-1","name":"Dungeon","place_id":987654}]}`),
	}}

	got, err := DiscoverStudioSessions(context.Background(), caller)
	if err != nil {
		t.Fatalf("DiscoverStudioSessions() error = %v", err)
	}
	if caller.name != "list_roblox_studios" {
		t.Fatalf("tool name = %q, want list_roblox_studios", caller.name)
	}
	if string(caller.arguments) != `{}` {
		t.Fatalf("arguments = %s, want {}", caller.arguments)
	}
	want := StudioSession{ID: "studio-1", Name: "Dungeon", PlaceID: "987654"}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("DiscoverStudioSessions() = %#v, want %#v", got, []StudioSession{want})
	}
}

func TestDiscoverStudioSessionsRejectsMCPErrorResult(t *testing.T) {
	caller := &fakeToolCaller{result: arenamcp.ToolResult{
		StructuredContent: json.RawMessage(`{"studios":[]}`),
		IsError:           true,
	}}

	_, err := DiscoverStudioSessions(context.Background(), caller)
	if err == nil {
		t.Fatal("DiscoverStudioSessions() error = nil, want MCP discovery error")
	}
}
