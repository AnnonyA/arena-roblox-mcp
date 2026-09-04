package agent

import "testing"

func TestContextKeepsNewestEventsWithinLimit(t *testing.T) {
	ctx := NewContext(2)
	ctx.Add(Event{Role: "user", Content: "one"})
	ctx.Add(Event{Role: "assistant", Content: "two"})
	ctx.Add(Event{Role: "tool", Content: "three"})

	got := ctx.Events()
	want := []Event{
		{Role: "assistant", Content: "two"},
		{Role: "tool", Content: "three"},
	}
	if len(got) != len(want) {
		t.Fatalf("Events() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Events()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}
