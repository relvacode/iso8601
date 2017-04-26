package iso8601

import (
	"encoding/json"
	"time"
)

var _ json.Unmarshaler = &Time{}

// Time is a helper object for parsing ISO8061 dates as a JSON string.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	} else {
		return ErrNotString
	}
	var err error
	t.Time, err = Parse(b)
	return err
}
