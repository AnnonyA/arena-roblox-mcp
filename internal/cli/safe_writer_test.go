package cli

import (
	"errors"
	"io"
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

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return len(p) - 1, nil
}

func TestWriteSafeMultilineReportsShortWrite(t *testing.T) {
	err := WriteSafeMultiline(shortWriter{}, "safe output")
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("WriteSafeMultiline() error = %v, want %v", err, io.ErrShortWrite)
	}
}
