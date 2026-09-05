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

func TestContextCompactsStaleToolOutputs(t *testing.T) {
	ctx := NewContext(4)
	ctx.Add(Event{Role: "user", Content: "keep user text"})
	ctx.Add(Event{Role: "tool", Content: "abcdefghij"})
	ctx.Add(Event{Role: "assistant", Content: "keep assistant text"})
	ctx.Add(Event{Role: "tool", Content: "latest tool output remains complete"})

	ctx.CompactToolOutputs(5)

	got := ctx.Events()
	want := []Event{
		{Role: "user", Content: "keep user text"},
		{Role: "tool", Content: "abcde… [compacted]"},
		{Role: "assistant", Content: "keep assistant text"},
		{Role: "tool", Content: "latest tool output remains complete"},
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

func TestContextCompactionPreservesUTF8(t *testing.T) {
	ctx := NewContext(2)
	ctx.Add(Event{Role: "tool", Content: "áéíóú"})
	ctx.Add(Event{Role: "assistant", Content: "done"})

	ctx.CompactToolOutputs(3)

	got := ctx.Events()
	if got[0].Content != "áéí… [compacted]" {
		t.Fatalf("compacted content = %q, want UTF-8-safe prefix", got[0].Content)
	}
}
