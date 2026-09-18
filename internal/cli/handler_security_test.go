package cli

import (
	"bytes"
	"testing"
)

func TestWriteSafeMultilinePreservesLayoutAndStripsTerminalControls(t *testing.T) {
	var out bytes.Buffer
	input := "arena.model=arena-code\nagent.safeMode=true\t# enabled\r\x1b[31mred\x1b[0m\u202Ehidden"

	if err := WriteSafeMultiline(&out, input); err != nil {
		t.Fatalf("WriteSafeMultiline returned error: %v", err)
	}

	const want = "arena.model=arena-code\nagent.safeMode=true\t# enabled[31mred[0mhidden"
	if got := out.String(); got != want {
		t.Fatalf("WriteSafeMultiline output = %q, want %q", got, want)
	}
}
