package cli

import (
	"reflect"
	"testing"
)

func TestSafeModelIDsDoesNotMutateInput(t *testing.T) {
	models := []string{"arena/escape\x1b[31m", "arena/safe-model"}
	original := append([]string(nil), models...)

	got := safeModelIDs(models)

	if !reflect.DeepEqual(models, original) {
		t.Fatalf("safeModelIDs() mutated input: got %q, want %q", models, original)
	}
	if want := []string{"arena/safe-model"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("safeModelIDs() = %q, want %q", got, want)
	}
}
