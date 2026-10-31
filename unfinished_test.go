package iso8601

import "testing"

func TestParseRejectsUnfinishedComponent(t *testing.T) {
	for _, prefix := range []string{
		"2024-01-01T12:", "2024-01-01T12:34:",
		"2024-01-01T12:34:56.", "2024-01-01T12:34:56,",
	} {
		for _, zone := range []string{"", "Z", "+01:00"} {
			input := prefix + zone
			t.Run(input, func(t *testing.T) {
				if got, err := ParseString(input); err == nil {
					t.Errorf("accepted unfinished input as %v", got)
				}
			})
		}
	}
}

func TestParseAcceptsCompleteZeroComponents(t *testing.T) {
	for _, input := range []string{
		"2024-01-01", "2024-01-01T", "2024-01-01 ", "2024-01-01TZ", "2024-01-01T00", "2024-01-01T00:00", "2024-01-01T00:00:00",
		"2024-01-01T00:00:00.0", "2024-01-01T00:00:00,0",
		"2024-01-01T00:00:00Z", "2024-01-01T00:00:00.0+01:00",
	} {
		t.Run(input, func(t *testing.T) {
			if _, err := ParseString(input); err != nil {
				t.Fatal(err)
			}
		})
	}
}
