package config

import (
	"strings"
	"testing"
)

func TestRejectDuplicateJSONKeysRejectsExcessiveNesting(t *testing.T) {
	data := []byte(`{"x":` + strings.Repeat("[", 65) + `0` + strings.Repeat("]", 65) + `}`)

	err := rejectDuplicateJSONKeys(data)
	if err == nil {
		t.Fatal("expected excessively nested JSON to be rejected")
	}
	if !strings.Contains(err.Error(), "nested too deeply") {
		t.Fatalf("expected nesting-limit error, got %v", err)
	}
}
