package cli

import (
	"strings"
	"testing"
)

func TestWriteSafeMultilineSanitizesTerminalControls(t *testing.T) {
	var out strings.Builder
	if err := WriteSafeMultiline(&out, "hello\x1b[31m\u200bworld\n\tcode"); err != nil {
		t.Fatalf("WriteSafeMultiline() error = %v", err)
	}
	const want = "hello[31mworld\n\tcode"
	if got := out.String(); got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestWriteSafeMultilineRejectsNilOutput(t *testing.T) {
	if err := WriteSafeMultiline(nil, "text"); err == nil {
		t.Fatal("WriteSafeMultiline() error = nil, want error")
	}
}
