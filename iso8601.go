// Package iso8061 is a utility for parsing ISO8601 datetime strings in Go.
// The standard library's RFC3339 reference layout is too strict for working with 3rd party APIs, especially ones written in other languages like Java.
//
// This library intends to support the full ISO8061 date specification with as much performance as possible.
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
	milliToNano uint = 1000000
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
//
// ParseISO zone will only look at the first 3 characters of the input string.
// The expectation is that all zone offsets are in hour intervals.
func ParseISOZone(inp string) (*time.Location, error) {
	if len(inp) < 3 {
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

	var z uint
	for i := 1; i < len(inp); i++ {
		if i == 3 {
			break
		}
		switch inp[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			z += uint(inp[i]) - charStart
		default:
			return nil, newUnexpectedCharacterError(inp[i])
		}
		z = z * 10
	}

	var offset int = int(z)
	if neg {
		offset = -offset
	}
	return time.FixedZone("", offset*360), nil
}

// Parse parses a full ISO8601 compliant date string into a time.Time object.
func Parse(inp string) (time.Time, error) {
	var (
		Y  uint
		M  uint
		d  uint
		h  uint
		m  uint
		s  uint
		ms uint
	)

	var loc = time.UTC

	var c uint
	var p uint = year

	var i int
parse:
	for ; i < len(inp); i++ {
		switch inp[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			c = c * 10
			c += uint(inp[i]) - charStart
		case '-':
			if p < second {
				switch p {
				case year:
					Y = c
				case month:
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
			case second:
				s = c
			case millisecond:
				ms = c
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
			case second:
				s = c
			case millisecond:
				ms = c
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
	// Sometimes a date can end without a `T` suffix.
	if c > 0 && p == day {
		d = c
	}
	return time.Date(int(Y), time.Month(M), int(d), int(h), int(m), int(s), int(ms*milliToNano), loc), nil
}
