package session

import "testing"

func TestResolveDataDirRejectsBlankUserConfigDirFallback(t *testing.T) {
	got, err := resolveDataDir("", func() (string, error) {
		return " \t ", nil
	})
	if err == nil {
		t.Fatalf("resolveDataDir() = %q, nil error; want blank fallback rejection", got)
	}
}
