package iso8601

import (
	"fmt"
	"testing"
	"time"
)

// ordinalCases covers ISO 8601 ordinal dates (YYYY-DDD), where the three digits
// after the year are the day of the year rather than a month.
var ordinalCases = []TestCase{
	{Using: "2020-001", Year: 2020, Month: 1, Day: 1},
	{Using: "2020-002", Year: 2020, Month: 1, Day: 2},
	{Using: "2020-012", Year: 2020, Month: 1, Day: 12},
	{Using: "2020-032", Year: 2020, Month: 2, Day: 1},
	{Using: "2020-060", Year: 2020, Month: 2, Day: 29}, // leap year, day 60
	{Using: "2020-061", Year: 2020, Month: 3, Day: 1},
	{Using: "2020-100", Year: 2020, Month: 4, Day: 9},
	{Using: "2020-366", Year: 2020, Month: 12, Day: 31}, // last day of a leap year
	{Using: "2019-060", Year: 2019, Month: 3, Day: 1},   // non-leap year, day 60
	{Using: "2019-365", Year: 2019, Month: 12, Day: 31},

	// Ordinal dates with a time and/or zone.
	{Using: "2020-012T10:20:30", Year: 2020, Month: 1, Day: 12, Hour: 10, Minute: 20, Second: 30},
	{Using: "2020-012T10:20:30Z", Year: 2020, Month: 1, Day: 12, Hour: 10, Minute: 20, Second: 30},
	{Using: "2020-012Z", Year: 2020, Month: 1, Day: 12},
	{Using: "2020-012+05:00", Year: 2020, Month: 1, Day: 12, Zone: 5},

	// Day of year out of range for the given year.
	{Using: "2020-000", ShouldInvalidRange: true, RangeElementWhenInvalid: "day"},
	{Using: "2020-367", ShouldInvalidRange: true, RangeElementWhenInvalid: "day"},
	{Using: "2021-366", ShouldInvalidRange: true, RangeElementWhenInvalid: "day"}, // 2021 is not a leap year
}

func TestOrdinal(t *testing.T) {
	for _, c := range ordinalCases {
		t.Run(c.Using, func(t *testing.T) {
			d, err := Parse([]byte(c.Using))
			if c.CheckError(err, t) {
				return
			}
			c.Check(d, t)
		})
	}
}

func TestOrdinalString(t *testing.T) {
	for _, c := range ordinalCases {
		t.Run(c.Using, func(t *testing.T) {
			d, err := ParseString(c.Using)
			if c.CheckError(err, t) {
				return
			}
			c.Check(d, t)
		})
	}
}

// malformedInputs are inputs that were previously accepted and returned a
// silently wrong time. Each must now fail to parse.
var malformedInputs = []string{
	"2020-01-02T16:20:45:99", // colon after the seconds field
	"2020-01-02T16:20:45:",   // trailing colon after the seconds field
	"2020-01-02T16::20",      // empty minute component
	"2020-01-02T16:20:.5",    // empty second component
	"2020--01-02",            // empty month component
}

func TestRejectMalformed(t *testing.T) {
	for _, s := range malformedInputs {
		t.Run(s, func(t *testing.T) {
			if d, err := ParseString(s); err == nil {
				t.Errorf("expected %q to fail parsing, got %s", s, d)
			}
		})
	}
}

// ordinalReference converts a day of the year to a calendar month and day by
// walking the months. It is deliberately independent of the parser's
// time.Date based normalisation so it can act as an oracle.
func ordinalReference(year, doy int) (month, day int) {
	month = 1
	for {
		dim := daysIn(time.Month(month), year)
		if doy <= dim {
			return month, doy
		}
		doy -= dim
		month++
	}
}

// TestOrdinalDayOfYear exhaustively checks every valid day of the year across
// leap and non-leap years (including century years) against the reference
// conversion, and checks that the day after the final day is out of range.
func TestOrdinalDayOfYear(t *testing.T) {
	for _, year := range []int{1900, 1999, 2000, 2019, 2020, 2021, 2024, 2100} {
		max := 365
		if isLeap(year) {
			max = 366
		}
		for doy := 1; doy <= max; doy++ {
			s := fmt.Sprintf("%04d-%03d", year, doy)
			got, err := ParseString(s)
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", s, err)
			}
			wantM, wantD := ordinalReference(year, doy)
			if got.Year() != year || int(got.Month()) != wantM || got.Day() != wantD {
				t.Errorf("%s = %04d-%02d-%02d; want %04d-%02d-%02d",
					s, got.Year(), int(got.Month()), got.Day(), year, wantM, wantD)
			}
		}
		s := fmt.Sprintf("%04d-%03d", year, max+1)
		if _, err := ParseString(s); err == nil {
			t.Errorf("%s: expected a day range error", s)
		}
	}
}
