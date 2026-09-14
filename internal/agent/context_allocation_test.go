package agent

import "testing"

func TestContextAddAtCapacityDoesNotAllocate(t *testing.T) {
	ctx := NewContext(2)
	ctx.Add(Event{Role: "user", Content: "one"})
	ctx.Add(Event{Role: "assistant", Content: "two"})

	allocs := testing.AllocsPerRun(100, func() {
		ctx.Add(Event{Role: "tool", Content: "steady-state"})
	})
	if allocs != 0 {
		t.Fatalf("Context.Add allocations at capacity = %.2f, want 0", allocs)
	}
}
