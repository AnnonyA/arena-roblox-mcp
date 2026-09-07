package arena

import (
	"net/http"
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

func TestRetryAfterDelayRespectsPastHTTPDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 7, 18, 0, 0, 0, time.UTC)
	past := now.Add(-time.Second).Format(http.TimeFormat)
	got := retryAfterDelay(past, now)
	if got != 0 {
		t.Fatalf("retry delay = %v, want 0 for past Retry-After HTTP-date", got)
	}
}

func TestRetryableStatusOnlyAccepts429And5xx(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		status int
		want   bool
	}{
		{status: http.StatusTooManyRequests, want: true},
		{status: http.StatusInternalServerError, want: true},
		{status: 599, want: true},
		{status: 600, want: false},
		{status: 700, want: false},
	} {
		if got := retryableStatus(tc.status); got != tc.want {
			t.Fatalf("retryableStatus(%d) = %v, want %v", tc.status, got, tc.want)
		}
	}
}
