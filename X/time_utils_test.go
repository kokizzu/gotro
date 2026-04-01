package X

import (
	"testing"
	"time"
)

func TestParseDateTimeVariants(t *testing.T) {
	loc := time.FixedZone("X", 7*3600)

	t.Run("date only", func(t *testing.T) {
		got, err := parseDateTime([]byte(`2024-03-01`), loc)
		if err != nil {
			t.Fatalf("parseDateTime failed: %v", err)
		}
		want := time.Date(2024, time.March, 1, 0, 0, 0, 0, loc)
		if !got.Equal(want) {
			t.Fatalf("date mismatch: want=%v got=%v", want, got)
		}
	})

	t.Run("datetime with 1-6 fractional digits", func(t *testing.T) {
		cases := map[string]int{
			`2024-03-01 02:03:04.1`:      100000000,
			`2024-03-01 02:03:04.12`:     120000000,
			`2024-03-01 02:03:04.123`:    123000000,
			`2024-03-01 02:03:04.1234`:   123400000,
			`2024-03-01 02:03:04.12345`:  123450000,
			`2024-03-01 02:03:04.123456`: 123456000,
		}
		for input, nsec := range cases {
			got, err := parseDateTime([]byte(input), loc)
			if err != nil {
				t.Fatalf("parseDateTime failed for %q: %v", input, err)
			}
			want := time.Date(2024, time.March, 1, 2, 3, 4, nsec, loc)
			if !got.Equal(want) {
				t.Fatalf("fraction mismatch for %q: want=%v got=%v", input, want, got)
			}
		}
	})

	t.Run("zero year month day fallback", func(t *testing.T) {
		got, err := parseDateTime([]byte(`0000-00-00 00:00:00.1`), loc)
		if err != nil {
			t.Fatalf("parseDateTime fallback failed: %v", err)
		}
		want := time.Date(1, time.January, 1, 0, 0, 0, 100000000, loc)
		if !got.Equal(want) {
			t.Fatalf("fallback mismatch: want=%v got=%v", want, got)
		}
	})

	t.Run("all-zero base returns zero time", func(t *testing.T) {
		got, err := parseDateTime([]byte(`0000-00-00 00:00:00.000000`), loc)
		if err != nil {
			t.Fatalf("parseDateTime failed: %v", err)
		}
		if !got.IsZero() {
			t.Fatalf("expected zero time, got %v", got)
		}
	})
}

func TestParseDateTimeErrors(t *testing.T) {
	loc := time.UTC
	cases := []string{
		`2024-03-01T00:00:00`,     // bad separator at pos 10
		`2024-03-01 00-00:00`,     // bad separator at pos 13
		`2024-03-01 00:00-00`,     // bad separator at pos 16
		`2024-03-01 00:00:00,123`, // bad separator at pos 19
	}
	for _, input := range cases {
		_, err := parseDateTime([]byte(input), loc)
		if err == nil {
			t.Fatalf("expected error for %q", input)
		}
	}
}
