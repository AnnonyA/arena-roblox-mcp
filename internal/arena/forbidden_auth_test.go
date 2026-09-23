package arena

import (
	"net/http"
	"testing"
)

func TestArenaStatusErrorTreatsForbiddenAsAuthenticationFailure(t *testing.T) {
	t.Parallel()

	err := arenaStatusError("chat", http.StatusForbidden)
	const want = "Arena authentication failed. Check your API key."
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err, want)
	}
}
