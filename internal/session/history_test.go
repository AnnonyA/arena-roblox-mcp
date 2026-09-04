package session

import "testing"

func TestHistoryKeepsOnlyNewestActionsWithinCapacity(t *testing.T) {
	history, err := NewHistory(2)
	if err != nil {
		t.Fatalf("NewHistory() error = %v", err)
	}

	history.Add(Action{Tool: "read_script", Summary: "read first"})
	history.Add(Action{Tool: "edit_script", Summary: "edit second"})
	history.Add(Action{Tool: "run_playtest", Summary: "test third"})

	got := history.Actions()
	want := []Action{
		{Tool: "edit_script", Summary: "edit second"},
		{Tool: "run_playtest", Summary: "test third"},
	}
	if len(got) != len(want) {
		t.Fatalf("len(Actions()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Actions()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}
