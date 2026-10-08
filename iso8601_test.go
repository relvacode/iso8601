package iso8601

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type TestCase struct {
	Using string

	Year  int
	Month int
	Day   int

	Hour        int
	Minute      int
	Second      int
	MilliSecond int

	Zone            float64
	ShouldFailParse bool

	ShouldInvalidRange      bool
	RangeElementWhenInvalid string
}

func (tc TestCase) CheckError(err error, t *testing.T) bool {
	if err != nil {
		if tc.ShouldInvalidRange {
			var re *RangeError
			if !errors.As(err, &re) {
				t.Fatalf("Found error %s of type %T but was expecting a RangeError", err, err)
			}

			if tc.RangeElementWhenInvalid != "" && re.Element != tc.RangeElementWhenInvalid {
				t.Fatalf("Expected a range error on %q but encountered %q: %s", tc.RangeElementWhenInvalid, re.Element, err)
			}

			return true
		}
		if tc.ShouldFailParse {
			return true
		}

		t.Fatal(err)
		return false
	}

	if err == nil && (tc.ShouldFailParse || tc.ShouldInvalidRange) {
		reason := "fail to parse"
		if tc.ShouldInvalidRange {
			reason = "to catch an invalid date range"
		}
		t.Fatalf("Expected test case %s", reason)
		return true
	}

	return false
}

func (tc TestCase) Check(d time.Time, t *testing.T) {
	if y := d.Year(); y != tc.Year {
		t.Errorf("Year = %d; want %d", y, tc.Year)
	}
	if m := int(d.Month()); m != tc.Month {
		t.Errorf("Month = %d; want %d", m, tc.Month)
	}
	if d := d.Day(); d != tc.Day {
		t.Errorf("Day = %d; want %d", d, tc.Day)
	}
	if h := d.Hour(); h != tc.Hour {
		t.Errorf("Hour = %d; want %d", h, tc.Hour)
	}
	if m := d.Minute(); m != tc.Minute {
		t.Errorf("Minute = %d; want %d", m, tc.Minute)
	}
	if s := d.Second(); s != tc.Second {
		t.Errorf("Second = %d; want %d", s, tc.Second)
	}

	if ms := d.Nanosecond() / 1000000; ms != tc.MilliSecond {
		t.Errorf(
			"Millisecond = %d; want %d (%d nanoseconds)",
			ms,
			tc.MilliSecond,
			d.Nanosecond(),
		)
	}

	_, z := d.Zone()
	if offset := float64(z) / 3600; offset != tc.Zone {
		t.Errorf("Zone = %.2f (%d); want %.2f", offset, z, tc.Zone)
	}
}

var cases = []TestCase{
	{
		Using: "2017-04-24T09:41:34.502+0100",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        1,
	},
	{
		Using: "2017-04-24T09:41+0100",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41,
		Zone: 1,
	},
	{
		Using: "2017-04-24T09+0100",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9,
		Zone: 1,
	},
	{
		Using: "2017-04-24T",
		Year:  2017, Month: 4, Day: 24,
	},
	{
		Using: "2017-04-24",
		Year:  2017, Month: 4, Day: 24,
	},
	{
		Using: "2017-04-24T09:41:34+0100",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		Zone: 1,
	},
	{
		Using: "2017-04-24T09:41:34.502-0100",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        -1,
	},
	{
		Using: "2017-04-24T09:41:34.502-01:00",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        -1,
	},
	{
		Using: "2017-04-24T09:41-01:00",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41,
		Zone: -1,
	},
	{
		Using: "2017-04-24T09-01:00",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9,
		Zone: -1,
	},
	{
		Using: "2017-04-24T09:41:34-0100",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		Zone: -1,
	},
	{
		Using: "2017-04-24T09:41:34.502Z",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        0,
	},
	{
		Using: "2017-04-24T09:41:34Z",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		Zone: 0,
	},
	{
		Using: "2017-04-24T09:41Z",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41,
		Zone: 0,
	},
	{
		Using: "2017-04-24T09Z",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9,
		Zone: 0,
	},
	{
		Using: "2017-04-24T09:41:34.089",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 89,
		Zone:        0,
	},
	{
		Using: "2017-04-24T09:41",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41,
		Zone: 0,
	},
	{
		Using: "2017-04-24T09",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9,
		Zone: 0,
	},
	{
		Using: "2017-04-24T09:41:34.009",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 9,
		Zone:        0,
	},
	{
		Using: "2017-04-24T09:41:34.893",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 893,
		Zone:        0,
	},
	{
		Using: "2017-04-24T09:41:34.89312523Z",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 893,
		Zone:        0,
	},
	{
		Using: "2017-04-24T09:41:34.502-0530",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        -5.5,
	},
	{
		Using: "2017-04-24T09:41:34.502+0530",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        5.5,
	},
	{
		Using: "2017-04-24T09:41:34.502+05:30",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        5.5,
	},

	{
		Using: "2017-04-24T09:41:34.502+05:45",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        5.75,
	},
	{
		Using: "2017-04-24T09:41:34.502+00",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        0,
	},
	{
		Using: "2017-04-24T09:41:34.502+0000",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        0,
	},
	{
		Using: "2017-04-24T09:41:34.502+00:00",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        0,
	},
	{
		Using: "+2017-04-24T09:41:34.502+00:00",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        0,
	},
	{
		Using: "2017",
		Year:  2017, Month: 1, Day: 1, Hour: 0, Minute: 0, Second: 0,
		MilliSecond: 0,
		Zone:        0,
	},
	{
		Using: "2017-02",
		Year:  2017, Month: 2, Day: 1, Hour: 0, Minute: 0, Second: 0,
		MilliSecond: 0,
		Zone:        0,
	},
	{
		Using: "2017-02-16",
		Year:  2017, Month: 2, Day: 16, Hour: 0, Minute: 0, Second: 0,
		MilliSecond: 0,
		Zone:        0,
	},
	{
		Using: "2017-04-24 09:41:34",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
	},
	{
		Using: "2017-04-24 09:41:34.502+00:00",
		Year:  2017, Month: 4, Day: 24,
		Hour: 9, Minute: 41, Second: 34,
		MilliSecond: 502,
		Zone:        0,
	},

	// Invalid Parse Test Cases
	{
		Using:           "2017-04-24T09:41:34.502-00",
		ShouldFailParse: true,
	},
	{
		Using:           "2017-04-24T09:41:34.502-0000",
		ShouldFailParse: true,
	},
	{
		Using:           "2017-04-24T09:41:34.502-00:00",
		ShouldFailParse: true,
	},
	{
		Using:           "-2017-04-24T09:41:34.502-00:00",
		ShouldFailParse: true,
	},
	{
		Using:           "2017-+04-24T09:41:34.502-00:00",
		ShouldFailParse: true,
	},
	{
		Using:           "2020-01-01T00:00:00+01:60",
		ShouldFailParse: true,
	},
	{
		Using:           "2020-01-01T00:00:00+24:00",
		ShouldFailParse: true,
	},
	{
		Using:           "2017-01-01T00:00:60.000Z+",
		ShouldFailParse: true,
	},
	{
		Using:           "2017-01-01T00:00:60.000Zz",
		ShouldFailParse: true,
	},
	{
		Using:           "2017-01-01T00:00:60.000Z00:00",
		ShouldFailParse: true,
	},

	// Invalid Range Test Cases
	{
		Using:                   "2017-00-01T00:00:00.000+00:00",
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "month",
	},
	{
		Using:                   "2017-13-01T00:00:00.000+00:00",
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "month",
	},

	{
		Using:                   "2017-01-00T00:00:00.000+00:00",
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "day",
	},
	{
		Using:                   "2017-01-32T00:00:00.000+00:00",
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "day",
	},
	{
		Using:                   "2019-02-29T00:00:00.000+00:00",
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "day",
	},
	{
		Using:                   "2020-02-30T00:00:00.000+00:00", // Leap year
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "day",
	},

	{
		Using:                   "2017-01-01T24:00:00.000+00:00",
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "hour",
	},

	{
		Using:                   "2017-01-01T00:60:00.000+00:00",
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "minute",
	},

	{
		Using:                   "2017-01-01T00:00:60.000+00:00",
		ShouldInvalidRange:      true,
		RangeElementWhenInvalid: "second",
	},
}

func TestParse(t *testing.T) {
	for _, c := range cases {
		t.Run(
			c.Using, func(t *testing.T) {
				d, err := Parse([]byte(c.Using))
				if c.CheckError(err, t) {
					return
				}
				t.Log(d)
				c.Check(d, t)
			},
		)

	}
}

func TestParseString(t *testing.T) {
	for _, c := range cases {
		t.Run(
			c.Using, func(t *testing.T) {
				d, err := ParseString(c.Using)
				if c.CheckError(err, t) {
					return
				}
				t.Log(d)
				c.Check(d, t)
			},
		)
	}
}

func BenchmarkParse(b *testing.B) {
	x := []byte("2017-04-24T09:41:34.502Z")
	for i := 0; i < b.N; i++ {
		_, err := Parse(x)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestParseStringInLocation(t *testing.T) {
	cases := []TestCase{
		{
			Using: "2017-04-24T09:41:34.502+05:45",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        5.75,
		},
		{
			Using: "2017-04-24T09:41:34.502",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        5,
		},
		{
			Using: "2017-04-24T09:41:34.502Z",
			Year:  2017, Month: 4, Day: 24,
			Hour: 9, Minute: 41, Second: 34,
			MilliSecond: 502,
			Zone:        0,
		},
	}

	loc := time.FixedZone("UTC+5", 5*60*60)

	for _, c := range cases {
		t.Run(
			c.Using, func(t *testing.T) {
				d, err := ParseStringInLocation(c.Using, loc)
				if c.CheckError(err, t) {
					return
				}
				t.Log(d)
				c.Check(d, t)
			},
		)
	}
}

type ZoneTestCase struct {
	Using  string
	Zone   float64
	Expect error
}

func TestParseISOZone(t *testing.T) {
	var zoneTestCases = []ZoneTestCase{
		{
			Using:  "",
			Expect: ErrZoneCharacters,
		},
		{
			Using: "Z",
			Zone:  0,
		},
		{
			Using: "z",
			Zone:  0,
		},
		{
			Using: "+00:00",
			Zone:  0,
		},
		{
			Using:  "00:00",
			Expect: UnexpectedCharacterError{Character: '0'},
		},
		{
			Using:  "-00:00",
			Expect: ErrInvalidZone,
		},
		{
			Using: "-05:30",
			Zone:  -5.5,
		},
		{
			Using: "+05:30",
			Zone:  5.5,
		},
		{
			Using: "-0530",
			Zone:  -5.5,
		},
		{
			Using:  "Zz",
			Expect: ErrRemainingData,
		},
		{
			Using:  "^",
			Expect: UnexpectedCharacterError{Character: '^'},
		},
		{
			Using: "-01",
			Zone:  -1,
		},
		{
			Using: "+23:00",
			Zone:  23,
		},
		{
			Using:  "+01:60",
			Expect: ErrInvalidZone,
		},
		{
			Using:  "+00:60",
			Expect: ErrInvalidZone,
		},
		{
			Using:  "+0160",
			Expect: ErrInvalidZone,
		},
		{
			Using:  "+24:00",
			Expect: ErrInvalidZone,
		},
	}

	for _, tc := range zoneTestCases {
		t.Run(tc.Using, func(t *testing.T) {
			z, err := ParseISOZone([]byte(tc.Using))
			if !errors.Is(err, tc.Expect) {
				t.Errorf("ParseISOZone expected to return error %v (%T), got %v (%T)", tc.Expect, tc.Expect, err, err)
				return
			}

			if tc.Expect != nil {
				return
			}

			ts := time.Date(2024, 1, 1, 1, 1, 1, 0, z)
			_, offset := ts.Zone()

			if offset := float64(offset) / 3600; offset != tc.Zone {
				t.Errorf("ParseISOZone expected to return zone %v, got %v", tc.Zone, offset)
			}
		})
	}
}

func TestParseISOZoneLayouts(t *testing.T) {
	invalid := []string{
		"+013", "-013",
		"+01:", "-01:",
		"+01:3", "-01:3",
		"+01003", "-01003",
		"+01000", "-01000",
		"+000", "-000",
		"+00:", "-00:",
		"+00:0", "-00:0",
	}
	for _, zone := range invalid {
		t.Run(zone, func(t *testing.T) {
			loc, err := ParseISOZone([]byte(zone))
			if err == nil || loc != nil {
				t.Errorf("ParseISOZone(%q) = %v, %v; want nil location and error", zone, loc, err)
			}

			input := "2024-01-02T03:04:05" + zone
			if got, err := ParseString(input); err == nil {
				t.Errorf("ParseString(%q) = %v; want error", input, got)
			}

			var got Time
			if err := json.Unmarshal([]byte(`"`+input+`"`), &got); err == nil {
				t.Errorf("UnmarshalJSON(%q) = %v; want error", input, got)
			}
		})
	}

	valid := []struct {
		zone   string
		offset int
	}{
		{"Z", 0},
		{"z", 0},
		{"+00", 0},
		{"+0000", 0},
		{"+00:00", 0},
		{"+01", 3600},
		{"-01", -3600},
		{"+0145", 6300},
		{"-0145", -6300},
		{"+01:45", 6300},
		{"-01:45", -6300},
		{"+23:59", 86340},
		{"-23:59", -86340},
	}
	for _, tc := range valid {
		t.Run(tc.zone, func(t *testing.T) {
			loc, err := ParseISOZone([]byte(tc.zone))
			if err != nil {
				t.Fatal(err)
			}
			_, offset := time.Date(2024, 1, 2, 3, 4, 5, 0, loc).Zone()
			if offset != tc.offset {
				t.Errorf("ParseISOZone(%q) offset = %d; want %d", tc.zone, offset, tc.offset)
			}
		})
	}
}

func TestParseFractionPrecision(t *testing.T) {
	valid := map[string]int{
		"2017-04-24T09:41:34.5Z":         500000000,
		"2017-04-24T09:41:34.000000001Z": 1,
		"2017-04-24T09:41:34.123456789Z": 123456789,
	}
	for in, want := range valid {
		d, err := ParseString(in)
		if err != nil {
			t.Fatalf("ParseString(%q): %v", in, err)
		}
		if d.Nanosecond() != want {
			t.Errorf("ParseString(%q) nanoseconds = %d, want %d", in, d.Nanosecond(), want)
		}
	}

	tooPrecise := []string{
		"2017-04-24T09:41:34.1234567891Z",
		// Leading zeros keep the accumulated value small, but the digits
		// still describe sub-nanosecond precision (0.1ns, not 1ns).
		"2017-04-24T09:41:34.0000000001Z",
		"2017-04-24T09:41:34.00000000012",
		// 2^64+1: overflows the accumulator back to 1.
		"2017-04-24T09:41:34.18446744073709551617Z",
	}
	for _, in := range tooPrecise {
		d, err := ParseString(in)
		if !errors.Is(err, ErrPrecision) {
			t.Errorf("ParseString(%q) = %v, %v; want ErrPrecision", in, d, err)
		}
	}
}
