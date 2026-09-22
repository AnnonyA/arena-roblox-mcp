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

func TestRetryAfterDelayTrimsHTTPDateWhitespace(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 7, 18, 0, 0, 0, time.UTC)
	future := now.Add(5 * time.Second).Format(http.TimeFormat)
	got := retryAfterDelay(" \t"+future+"\t ", now)
	if got != 5*time.Second {
		t.Fatalf("retry delay = %v, want %v for whitespace-padded Retry-After HTTP-date", got, 5*time.Second)
	}
}

func TestRetryAfterDelayBoundsFarFutureHTTPDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.September, 7, 18, 0, 0, 0, time.UTC)
	future := now.Add(24 * time.Hour).Format(http.TimeFormat)
	got := retryAfterDelay(future, now)
	if got != time.Minute {
		t.Fatalf("retry delay = %v, want maximum %v for far-future Retry-After HTTP-date", got, time.Minute)
	}
}

func TestRetryAfterDelayBoundsLargeNumericDelay(t *testing.T) {
	t.Parallel()

	got := retryAfterDelay("3600", time.Unix(0, 0))
	if got != time.Minute {
		t.Fatalf("retry delay = %v, want maximum %v for large numeric Retry-After", got, time.Minute)
	}
}

func TestRetryAfterDelayCapsDurationOverflow(t *testing.T) {
	t.Parallel()

	got := retryAfterDelay("9223372036854775807", time.Unix(0, 0))
	if got != maxRetryDelay {
		t.Fatalf("retry delay = %v, want maximum %v for overflowing Retry-After", got, maxRetryDelay)
	}
}

func TestRetryAfterDelayCapsIntegerParseOverflow(t *testing.T) {
	t.Parallel()

	got := retryAfterDelay("92233720368547758070", time.Unix(0, 0))
	if got != maxRetryDelay {
		t.Fatalf("retry delay = %v, want maximum %v for integer-overflowing Retry-After", got, maxRetryDelay)
	}
}

func FuzzRetryAfterDelayStaysBounded(f *testing.F) {
	for _, seed := range []string{"", "0", "2", "3600", "9223372036854775807", "Wed, 21 Oct 2015 07:28:00 GMT", "not-a-delay"} {
		f.Add(seed)
	}

	now := time.Date(2026, time.September, 7, 18, 0, 0, 0, time.UTC)
	f.Fuzz(func(t *testing.T, value string) {
		got := retryAfterDelay(value, now)
		if got < 0 || got > maxRetryDelay {
			t.Fatalf("retryAfterDelay(%q) = %v, want delay in [0, %v]", value, got, maxRetryDelay)
		}
	})
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
