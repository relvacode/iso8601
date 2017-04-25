package iso8601

import (
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
