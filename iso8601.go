// Package iso8601 is a utility for parsing ISO8601 datetime strings into native Go times.
// The standard library's RFC3339 reference layout can be too strict for working with 3rd party APIs,
// especially ones written in other languages.
//
// Use the provided `Time` structure instead of the default `time.Time` to provide ISO8601 support for JSON responses.
//
package iso8601

import (
	"time"
)

const (
	year uint = iota
	month
	day
	hour
	minute
	second
	millisecond
)

const (
	// charStart is the binary position of the character `0`
	charStart uint = 48
)

// ParseISOZone parses the 5 character zone information in an ISO8061 date string.
// This function expects input that matches:
//
//     -0100
//     +0100
//     +01:00
//     -01:00
//     +01
//     +01:45
//     +0145
//
func ParseISOZone(inp []byte) (*time.Location, error) {
	if len(inp) < 3 || len(inp) > 6 {
		return nil, ErrZoneCharacters
	}
	var neg bool
	switch inp[0] {
	case '+':
	case '-':
		neg = true
	default:
		return nil, newUnexpectedCharacterError(inp[0])
	}

	var offset int

	var z uint
	var multiplier = uint(3600) // start with initial multiplier of hours
	for i := 1; i < len(inp); i++ {
		if i == 3 { // next multiplier
			offset = int(z * multiplier)
			multiplier = 60 // multiplier for minutes
			z = 0
		} else { // next digit
			z = z * 10
		}

		switch inp[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			z += uint(inp[i]) - charStart
		case ':':
			if i != 3 {
				return nil, newUnexpectedCharacterError(inp[i])
			}
		default:
			return nil, newUnexpectedCharacterError(inp[i])
		}

	}

	offset += int(z * multiplier)

	if neg {
		offset = -offset
	}
	return time.FixedZone("", offset), nil
}

// Parse parses an ISO8601 compliant date-time byte slice into a time.Time object.
func Parse(inp []byte) (time.Time, error) {
	var (
		Y         uint
		M         uint
		d         uint
		h         uint
		m         uint
		s         uint
		fraction  int
		nfraction = 1 //counts amount of precision for the second fraction
	)

	// Always assume UTC by default
	var loc = time.UTC

	var c uint
	var p = year

	var i int

	var lastnum uint
parse:
	for ; i < len(inp); i++ {
		switch inp[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			c = c * 10
			if p == day && uint(inp[i]) == 48 {
				if lastnum == 48 {
					c = 1
				}
				lastnum = 48
			} else {
				c += uint(inp[i]) - charStart
			}

			if p == millisecond {
				nfraction++
			}
		case '-':
			if p < hour {
				switch p {
				case year:
					if c == 0 {
						c = 1
					}
					Y = c
				case month:
					if c == 0 {
						c = 1
					}
					M = c
				default:
					return time.Time{}, newUnexpectedCharacterError(inp[i])
				}
				p++
				c = 0
				continue
			}
			fallthrough
		case '+':
			switch p {
			case hour:
				h = c
			case minute:
				m = c
			case second:
				s = c
			case millisecond:
				fraction = int(c)
			default:
				return time.Time{}, newUnexpectedCharacterError(inp[i])
			}
			c = 0
			var err error
			loc, err = ParseISOZone(inp[i:])
			if err != nil {
				return time.Time{}, err
			}
			break parse
		case 'T':
			if p != day {
				return time.Time{}, newUnexpectedCharacterError(inp[i])
			}
			d = c
			c = 0
			p++
		case ':':
			switch p {
			case hour:
				h = c
			case minute:
				m = c
			case second:
				m = c
			default:
				return time.Time{}, newUnexpectedCharacterError(inp[i])
			}
			c = 0
			p++
		case '.':
			if p != second {
				return time.Time{}, newUnexpectedCharacterError(inp[i])
			}
			s = c
			c = 0
			p++
		case 'Z':
			switch p {
			case hour:
				h = c
			case minute:
				m = c
			case second:
				s = c
			case millisecond:
				fraction = int(c)
			default:
				return time.Time{}, newUnexpectedCharacterError(inp[i])
			}
			c = 0
			if len(inp) != i+1 {
				return time.Time{}, ErrRemainingData
			}
		default:
			return time.Time{}, newUnexpectedCharacterError(inp[i])
		}
	}

	// Capture remaining data
	// Sometimes a date can end without a non-integer character
	if c > 0 {
		switch p {
		case day:
			d = c
		case hour:
			h = c
		case minute:
			m = c
		case second:
			s = c
		case millisecond:
			fraction = int(c)
		}
	}

	// Get the seconds fraction as nanoseconds
	if fraction < 0 || 1e9 <= fraction {
		return time.Time{}, ErrPrecision
	}
	scale := 10 - nfraction
	for i := 0; i < scale; i++ {
		fraction *= 10
	}
	return time.Date(int(Y), time.Month(M), int(d), int(h), int(m), int(s), fraction, loc), nil
}

// ParseString parses an ISO8601 compliant date-time string into a time.Time object.
func ParseString(inp string) (time.Time, error) {
	return Parse([]byte(inp))
}
