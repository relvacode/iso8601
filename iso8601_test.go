package iso8601

import (
	"testing"
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

	Zone int
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
}

func Assert(t *testing.T, x, y int) {
	if x != y {
		t.Fatalf("Expected %d but got %d", x, y)
	}
}

func TestParse(t *testing.T) {
	for _, c := range cases {
		t.Run(c.Using, func(t *testing.T) {
			d, err := Parse([]byte(c.Using))
			if err != nil {
				t.Fatal(err)
			}
			t.Log(d)
			Assert(t, c.Year, d.Year())
			Assert(t, c.Month, int(d.Month()))
			Assert(t, c.Day, d.Day())
			Assert(t, c.Hour, d.Hour())
			Assert(t, c.Minute, d.Minute())
			Assert(t, c.Second, d.Second())
			Assert(t, c.MilliSecond, d.Nanosecond()/int(milliToNano))

			_, z := d.Zone()
			Assert(t, c.Zone, z/3600)
		})

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
