package iso8601

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"
)

type TestAPIResponse struct {
	Ptr  *Time
	Nptr Time
}

type TestStdLibAPIResponse struct {
	Ptr  *time.Time
	Nptr time.Time
}

var ShortTest = TestCase{
	Using: "2001-11-13",
	Year:  2001, Month: 11, Day: 13,
}

var StructTestData = []byte(`
{
  "Ptr": "2017-04-26T11:13:04+01:00",
  "Nptr": "2017-04-26T11:13:04+01:00"
}
`)

var NullTestData = []byte(`
{
  "Ptr": null,
  "Nptr": null
}
`)

var ZeroedTestData = []byte(`
{
  "Ptr": "0001-01-01",
  "Nptr": "0001-01-01"
}
`)

var StructTest = TestCase{
	Year: 2017, Month: 04, Day: 26,
	Hour: 11, Minute: 13, Second: 04,
	Zone: 1,
}

func TestTime_Unmarshaling(t *testing.T) {
	t.Run("short", func(t *testing.T) {
		var b = []byte(`"2001-11-13"`)

		tn := new(Time)
		if err := tn.UnmarshalJSON(b); err != nil {
			t.Fatal(err)
		}

		if y := tn.Year(); y != ShortTest.Year {
			t.Errorf("Year = %d; want %d", y, ShortTest.Year)
		}

		if m := int(tn.Month()); m != ShortTest.Month {
			t.Errorf("Month = %d; want %d", m, ShortTest.Month)
		}

		if d := tn.Day(); d != ShortTest.Day {
			t.Errorf("Day = %d; want %d", d, ShortTest.Day)
		}

		err := tn.UnmarshalJSON([]byte(`2001-11-13`))
		if err != ErrNotString {
			t.Fatal(err)
		}
		if err == nil {
			t.Fatal("Expected an error from unmarshal")
		}
	})

	t.Run("struct", func(t *testing.T) {
		resp := new(TestAPIResponse)
		if err := json.Unmarshal(StructTestData, resp); err != nil {
			t.Fatal(err)
		}

		stdlibResp := new(TestStdLibAPIResponse)
		if err := json.Unmarshal(StructTestData, stdlibResp); err != nil {
			t.Fatal(err)
		}

		t.Run("stblib parity", func(t *testing.T) {
			if !resp.Ptr.Time.Equal(*stdlibResp.Ptr) || !resp.Nptr.Time.Equal(stdlibResp.Nptr) {
				t.Fatalf("Parsed time values are not equal to standard library implementation")
			}
		})

		t.Run("ptr", func(t *testing.T) {
			if y := resp.Ptr.Year(); y != StructTest.Year {
				t.Errorf("Ptr: Year = %d; want %d", y, StructTest.Year)
			}
			if d := resp.Ptr.Day(); d != StructTest.Day {
				t.Errorf("Ptr: Day = %d; want %d", d, StructTest.Day)
			}
			if s := resp.Ptr.Second(); s != StructTest.Second {
				t.Errorf("Ptr: Second = %d; want %d", s, StructTest.Second)
			}
		})

		t.Run("noptr", func(t *testing.T) {
			if y := resp.Nptr.Year(); y != StructTest.Year {
				t.Errorf("NoPtr: Year = %d; want %d", y, StructTest.Year)
			}
			if d := resp.Nptr.Day(); d != StructTest.Day {
				t.Errorf("NoPtr: Day = %d; want %d", d, StructTest.Day)
			}
			if s := resp.Nptr.Second(); s != StructTest.Second {
				t.Errorf("NoPtr: Second = %d; want %d", s, StructTest.Second)
			}
		})
	})

	t.Run("null", func(t *testing.T) {
		resp := new(TestAPIResponse)
		if err := json.Unmarshal(NullTestData, resp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("time zeroed", func(t *testing.T) {
		resp := new(TestAPIResponse)
		if err := json.Unmarshal(ZeroedTestData, resp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("reparse", func(t *testing.T) {
		s := time.Now().UTC()
		data := []byte(s.Format(time.RFC3339Nano))
		n, err := Parse(data)
		if err != nil {
			t.Fatal(err)
		}
		if !s.Equal(n) {
			t.Fatalf("Parsing a JSON date mismatch; wanted %s; got %s", s, n)
		}
	})

	t.Run("string", func(t *testing.T) {
		t1 := time.Now().UTC()
		s := Time{Time: t1}.String()
		expected := t1.Format(time.RFC3339Nano)
		if s != expected {
			t.Fatalf("String; wanted %s; got %s", expected, s)
		}
	})
}

func TestTime_Marshaling(t *testing.T) {
	t9 := Date(2017, 4, 26, 11, 13, 4, 123456789, time.UTC)

	cases := []struct {
		format     string
		resolution time.Duration
		expected   string
	}{
		{
			format:     RFC3339,
			resolution: time.Second,
			expected:   "2017-04-26T11:13:04Z",
		},
		{
			format:     RFC3339Milli,
			resolution: time.Millisecond,
			expected:   "2017-04-26T11:13:04.123Z",
		},
		{
			format:     RFC3339Micro,
			resolution: time.Microsecond,
			expected:   "2017-04-26T11:13:04.123456Z",
		},
		{
			format:     RFC3339Nano,
			resolution: time.Nanosecond,
			expected:   "2017-04-26T11:13:04.123456789Z",
		},
	}

	t.Run("text marshal/unmarshal", func(t *testing.T) {
		for _, c := range cases {
			MarshalTextFormat = c.format

			b, err := xml.Marshal(t9)
			if err != nil {
				t.Fatal(err)
			}

			expectedXML := fmt.Sprintf("<Time>%s</Time>", c.expected)
			if string(b) != expectedXML {
				t.Fatalf("wanted %s; got %s", expectedXML, string(b))
			}

			tn := Time{}
			if err := xml.Unmarshal(b, &tn); err != nil {
				t.Fatal(err)
			}

			expectedTime := t9.Truncate(c.resolution)

			if !tn.Time.Equal(expectedTime.Time) {
				t.Fatalf("wanted %s; got %s", expectedTime.Time, tn)
			}
		}
	})

	t.Run("JSON marshal/unmarshal", func(t *testing.T) {
		for _, c := range cases {
			MarshalTextFormat = c.format

			b, err := json.Marshal(t9)
			if err != nil {
				t.Fatal(err)
			}

			expectedJSON := fmt.Sprintf("%q", c.expected)
			if string(b) != expectedJSON {
				t.Fatalf("wanted %s; got %s", expectedJSON, string(b))
			}

			tn := new(Time)
			if err := tn.UnmarshalJSON(b); err != nil {
				t.Fatal(err)
			}

			expectedTime := t9.Truncate(c.resolution)

			if !tn.Time.Equal(expectedTime.Time) {
				t.Fatalf("wanted %s; got %s", expectedTime.Time, tn)
			}
		}
	})

	t.Run("Truncate", func(t *testing.T) {
		tr := t9.Truncate(time.Microsecond)

		expected := t9.Time.Truncate(time.Microsecond)

		if !tr.Time.Equal(expected) {
			t.Fatalf("wanted %s; got %s", expected, tr)
		}
	})

	t.Run("Round", func(t *testing.T) {
		r := t9.Round(time.Microsecond)

		expected := t9.Time.Round(time.Microsecond)

		if !r.Time.Equal(expected) {
			t.Fatalf("wanted %s; got %s", expected, r)
		}
	})
}

func BenchmarkCheckNull(b *testing.B) {
	var n = []byte("null")

	b.Run("compare", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bytes.Compare(n, n)
		}
	})
	b.Run("exact", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			null(n)
		}
	})
}
