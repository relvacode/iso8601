package iso8601

import (
	"encoding/json"
	"testing"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	var b = []byte(`"2001-11-13"`)

	tn := new(Time)
	if err := tn.UnmarshalJSON(b); err != nil {
		t.Fatal(err)
	}
	Assert(t, 2001, tn.Year())
	Assert(t, 11, int(tn.Month()))
	Assert(t, 13, tn.Day())

	err := tn.UnmarshalJSON([]byte(`2001-11-13`))
	if err != ErrNotString {
		t.Fatal(err)
	}
	if err == nil {
		t.Fatal("Expected an error")
	}
}

type TestAPIResponse struct {
	Ptr  *Time
	Nptr Time
}

var TestAPIData = []byte(`
{
  "Ptr": "2017-04-26T11:13:04+01:00",
  "Nptr": "2017-04-26T11:13:04+01:00"
}
`)

func TestTime_UnmarshalJSON2(t *testing.T) {
	resp := new(TestAPIResponse)
	if err := json.Unmarshal(TestAPIData, resp); err != nil {
		t.Fatal(err)
	}
	Assert(t, 2017, resp.Ptr.Year())
	Assert(t, 26, resp.Ptr.Day())
	Assert(t, 04, resp.Ptr.Second())

	Assert(t, 2017, resp.Nptr.Year())
	Assert(t, 26, resp.Nptr.Day())
	Assert(t, 04, resp.Nptr.Second())
}
