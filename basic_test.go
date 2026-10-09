package iso8601

import (
	"encoding/json"
	"testing"
	"time"
)

var basicCases = []TestCase{
	// Basic calendar date only (YYYYMMDD).
	{Using: "20170424", Year: 2017, Month: 4, Day: 24},
	{Using: "20200102", Year: 2020, Month: 1, Day: 2},
	{Using: "19991231", Year: 1999, Month: 12, Day: 31},
	{Using: "00010101", Year: 1, Month: 1, Day: 1},

	// Leap year verification.
	{Using: "20200229", Year: 2020, Month: 2, Day: 29},
	{Using: "20240229", Year: 2024, Month: 2, Day: 29},
	{Using: "20000229", Year: 2000, Month: 2, Day: 29},
	{Using: "19000228", Year: 1900, Month: 2, Day: 28},

	// Basic calendar date with leading plus.
	{Using: "+20170424", Year: 2017, Month: 4, Day: 24},

	// Basic calendar date with trailing delimiter and no time.
	{Using: "20170424T", Year: 2017, Month: 4, Day: 24},
	{Using: "20170424 ", Year: 2017, Month: 4, Day: 24},

	// Basic calendar date with time.
	{
		Using:  "20170424T09:41:34",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
		Second: 34,
	},
	{
		Using:  "20170424 09:41:34",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
		Second: 34,
	},
	{
		Using:  "20170424T09:41",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
	},
	{
		Using: "20170424T09",
		Year:  2017,
		Month: 4,
		Day:   24,
		Hour:  9,
	},

	// Basic calendar date with time and UTC zone.
	{
		Using:  "20170424T09:41:34Z",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
		Second: 34,
		Zone:   0,
	},
	{
		Using:  "20170424T09:41Z",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
		Zone:   0,
	},
	{
		Using: "20170424T09Z",
		Year:  2017,
		Month: 4,
		Day:   24,
		Hour:  9,
		Zone:  0,
	},

	// Basic calendar date with time and zone offsets.
	{
		Using:  "20170424T09:41:34+0100",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
		Second: 34,
		Zone:   1,
	},
	{
		Using:  "20170424T09:41:34+01:00",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
		Second: 34,
		Zone:   1,
	},
	{
		Using:  "20170424T09:41:34-0100",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
		Second: 34,
		Zone:   -1,
	},
	{
		Using:  "20170424T09:41:34-01:00",
		Year:   2017,
		Month:  4,
		Day:    24,
		Hour:   9,
		Minute: 41,
		Second: 34,
		Zone:   -1,
	},

	// Basic calendar date with fractional seconds.
	{
		Using:       "20170424T09:41:34.502",
		Year:        2017,
		Month:       4,
		Day:         24,
		Hour:        9,
		Minute:      41,
		Second:      34,
		MilliSecond: 502,
		Zone:        0,
	},
	{
		Using:       "20170424T09:41:34.502Z",
		Year:        2017,
		Month:       4,
		Day:         24,
		Hour:        9,
		Minute:      41,
		Second:      34,
		MilliSecond: 502,
		Zone:        0,
	},
	{
		Using:       "20170424T09:41:34.502+0100",
		Year:        2017,
		Month:       4,
		Day:         24,
		Hour:        9,
		Minute:      41,
		Second:      34,
		MilliSecond: 502,
		Zone:        1,
	},
	{
		Using:       "20170424T09:41:34.502+05:45",
		Year:        2017,
		Month:       4,
		Day:         24,
		Hour:        9,
		Minute:      41,
		Second:      34,
		MilliSecond: 502,
		Zone:        5.75,
	},
	{
		Using:       "20170424T09:41:34.502-0530",
		Year:        2017,
		Month:       4,
		Day:         24,
		Hour:        9,
		Minute:      41,
		Second:      34,
		MilliSecond: 502,
		Zone:        -5.5,
	},
	{
		Using:       "+20170424T09:41:34.502+00:00",
		Year:        2017,
		Month:       4,
		Day:         24,
		Hour:        9,
		Minute:      41,
		Second:      34,
		MilliSecond: 502,
		Zone:        0,
	},

	// Basic calendar date with zone only.
	{Using: "20170424Z", Year: 2017, Month: 4, Day: 24, Zone: 0},
	{Using: "20170424+05:00", Year: 2017, Month: 4, Day: 24, Zone: 5},
	{Using: "20170424-0500", Year: 2017, Month: 4, Day: 24, Zone: -5},

	// Out of range calendar components.
	{Using: "20170024", ShouldInvalidRange: true, RangeElementWhenInvalid: "month"},
	{Using: "20171324", ShouldInvalidRange: true, RangeElementWhenInvalid: "month"},
	{Using: "20170100", ShouldInvalidRange: true, RangeElementWhenInvalid: "day"},
	{Using: "20170132", ShouldInvalidRange: true, RangeElementWhenInvalid: "day"},
	{Using: "20190229", ShouldInvalidRange: true, RangeElementWhenInvalid: "day"},
	{Using: "19000229", ShouldInvalidRange: true, RangeElementWhenInvalid: "day"},
	{Using: "20170431", ShouldInvalidRange: true, RangeElementWhenInvalid: "day"},
	{Using: "20170424T24:00:00", ShouldInvalidRange: true, RangeElementWhenInvalid: "hour"},
	{Using: "20170424T00:60:00", ShouldInvalidRange: true, RangeElementWhenInvalid: "minute"},
	{Using: "20170424T00:00:60", ShouldInvalidRange: true, RangeElementWhenInvalid: "second"},
}

func TestBasic(t *testing.T) {
	for _, c := range basicCases {
		t.Run(c.Using, func(t *testing.T) {
			d, err := Parse([]byte(c.Using))
			if c.CheckError(err, t) {
				return
			}
			c.Check(d, t)
		})
	}
}

func TestBasicString(t *testing.T) {
	for _, c := range basicCases {
		t.Run(c.Using, func(t *testing.T) {
			d, err := ParseString(c.Using)
			if c.CheckError(err, t) {
				return
			}
			c.Check(d, t)
		})
	}
}

func TestBasicStringInLocation(t *testing.T) {
	loc := time.FixedZone("UTC+5", 5*60*60)

	cases := []TestCase{
		{Using: "20170424", Year: 2017, Month: 4, Day: 24, Zone: 5},
		{
			Using:  "20170424T09:41:34",
			Year:   2017,
			Month:  4,
			Day:    24,
			Hour:   9,
			Minute: 41,
			Second: 34,
			Zone:   5,
		},
		{
			Using:  "20170424T09:41:34Z",
			Year:   2017,
			Month:  4,
			Day:    24,
			Hour:   9,
			Minute: 41,
			Second: 34,
			Zone:   0,
		},
	}

	for _, c := range cases {
		t.Run(c.Using, func(t *testing.T) {
			d, err := ParseStringInLocation(c.Using, loc)
			if c.CheckError(err, t) {
				return
			}
			c.Check(d, t)
		})
	}
}

var basicMalformedInputs = []string{
	// ISO 8601 forbids basic year-month representation YYYYMM without hyphens.
	"201704",
	"201704T",
	"201704T09:41:34",
	"201704241T09:41:34",
	"20170424T09:41:34:00",
	"20170424T09:41:34:",
}

func TestBasicRejectMalformed(t *testing.T) {
	for _, s := range basicMalformedInputs {
		t.Run(s, func(t *testing.T) {
			if d, err := ParseString(s); err == nil {
				t.Errorf("expected %q to fail parsing, got %s", s, d)
			}
		})
	}
}

func TestBasicJSON(t *testing.T) {
	cases := []struct {
		input string
		want  time.Time
	}{
		{
			input: `"20170424"`,
			want:  time.Date(2017, 4, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			input: `"20170424T09:41:34Z"`,
			want:  time.Date(2017, 4, 24, 9, 41, 34, 0, time.UTC),
		},
		{
			input: `"20170424T09:41:34+01:00"`,
			want:  time.Date(2017, 4, 24, 9, 41, 34, 0, time.FixedZone("", 3600)),
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			var val Time
			if err := json.Unmarshal([]byte(tc.input), &val); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v", tc.input, err)
			}
			if !val.Equal(tc.want) {
				t.Errorf("Unmarshal(%s) = %v, want %v", tc.input, val.Time, tc.want)
			}
		})
	}
}

func BenchmarkParseBasic(b *testing.B) {
	x := []byte("20170424T09:41:34.502Z")
	for i := 0; i < b.N; i++ {
		_, err := Parse(x)
		if err != nil {
			b.Fatal(err)
		}
	}
}
