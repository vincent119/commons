package timex

import (
	"regexp"
	"testing"
	"time"
)

func TestNowUTC(t *testing.T) {
	now := NowUTC()
	if now.Location() != time.UTC {
		t.Error("NowUTC() should return time in UTC")
	}
}

func TestStartOfDay(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	// 2025-08-19 10:00:00+08 測試：零點應回到同一天 00:00:00+08（再轉 UTC）
	in := time.Date(2025, 8, 19, 10, 0, 0, 0, loc)
	got := StartOfDay(in, loc)

	// 期望的當地零點（轉 UTC）
	// 2025-08-19 00:00:00+08 = 2025-08-18 16:00:00 UTC
	want := time.Date(2025, 8, 18, 16, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("StartOfDay() got %v; want %v", got, want)
	}
}

func TestTruncateTo(t *testing.T) {
	// 2025-08-19 10:30:45+08
	in := time.Date(2025, 8, 19, 10, 30, 45, 0, time.UTC)

	// Truncate to hour -> 10:00:00
	want := time.Date(2025, 8, 19, 10, 0, 0, 0, time.UTC)
	got := TruncateTo(in, time.Hour)

	if !got.Equal(want) {
		t.Errorf("TruncateTo(Hour) got %v; want %v", got, want)
	}

	if got.Location() != time.UTC {
		t.Error("TruncateTo output should be UTC")
	}
}

func TestFormatTime(t *testing.T) {
	in := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	layout := "2006-01-02"
	want := "2025-01-01"

	if got := FormatTime(in, layout); got != want {
		t.Errorf("FormatTime(%v, %q) = %q, want %q", in, layout, got, want)
	}
}

func TestParseTime(t *testing.T) {
	str := "2025-01-01"
	layout := "2006-01-02"
	want := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	got, err := ParseTime(str, layout)
	if err != nil {
		t.Fatalf("ParseTime error: %v", err)
	}
	if !got.Equal(want) {
		t.Errorf("ParseTime got %v, want %v", got, want)
	}

	_, err = ParseTime("invalid", layout)
	if err == nil {
		t.Error("ParseTime should fail for invalid input")
	}
}

// Helper to match regex
func assertMatch(t *testing.T, str, pattern string) {
	t.Helper()
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		t.Fatalf("Regex error: %v", err)
	}
	if !matched {
		t.Errorf("String %q did not match pattern %q", str, pattern)
	}
}

func TestTimeStamp(t *testing.T) {
	// 2006-01-02T15:04:05.000Z07:00
	// Regex: ^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}[+-]\d{2}:?\d{2}$ (rough)
	// Output looks like: 2026-01-09T01:21:55.772+08:00
	assertMatch(t, TimeStamp(), `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}(Z|[+-]\d{2}:?\d{2})$`)
}

func TestTimeStampUTC(t *testing.T) {
	// 2006-01-02T15:04:05.000Z
	assertMatch(t, TimeStampUTC(), `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`)
}

func TestDateStamp(t *testing.T) {
	// 2006-01-02
	assertMatch(t, DateStamp(), `^\d{4}-\d{2}-\d{2}$`)
}

func TestUnixTimeStamp(t *testing.T) {
	assertMatch(t, UnixTimeStamp(), `^\d+$`)
}

func TestUnixMilliStamp(t *testing.T) {
	assertMatch(t, UnixMilliStamp(), `^\d+$`)
}

func TestTimeOnlyStamp(t *testing.T) {
	// 15:04:05
	assertMatch(t, TimeOnlyStamp(), `^\d{2}:\d{2}:\d{2}$`)
}

func TestWithZoneTimeStamp(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	// 2006-01-02T15:04:05.000Z07:00
	ts := WithZoneTimeStamp(loc)
	assertMatch(t, ts, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+08:00$`)
}

func TestFormatISO8601(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	// 2026-01-12 18:09:11 +0800
	in := time.Date(2026, 1, 12, 18, 9, 11, 0, loc)

	got := FormatISO8601(in)
	want := "2026-01-12T18:09:11.000+0800"

	if got != want {
		t.Errorf("FormatISO8601() = %q, want %q", got, want)
	}

	// Regex check for general validity
	assertMatch(t, got, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}[+-]\d{4}$`)
}
