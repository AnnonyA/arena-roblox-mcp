package arena

import (
	"testing"
	"time"
)

func TestRetryAfterDelayUsesFallbackWhenHeaderMissing(t *testing.T) {
	t.Parallel()

	got := retryAfterDelay("", time.Unix(0, 0))
	if got <= 0 {
		t.Fatalf("retry delay = %v, want positive fallback", got)
	}
}

func TestRetryAfterDelayPrefersServerValue(t *testing.T) {
	t.Parallel()

	got := retryAfterDelay("2", time.Unix(0, 0))
	if got != 2*time.Second {
		t.Fatalf("retry delay = %v, want %v", got, 2*time.Second)
	}
}

func TestRetryAfterDelayRespectsExplicitZero(t *testing.T) {
	t.Parallel()

	got := retryAfterDelay("0", time.Unix(0, 0))
	if got != 0 {
		t.Fatalf("retry delay = %v, want 0 for explicit Retry-After zero", got)
	}
}
