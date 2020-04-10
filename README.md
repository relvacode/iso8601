# iso8601
A fast ISO8601 date parser for Go

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/relvacode/iso8601) [![GoDoc](https://godoc.org/github.com/relvacode/iso8601?status.svg)](https://godoc.org/github.com/relvacode/iso8601) [![Build Status](https://travis-ci.org/relvacode/iso8601.svg?branch=master)](https://travis-ci.org/relvacode/iso8601) [![Go Report Card](https://goreportcard.com/badge/github.com/relvacode/iso8601)](https://goreportcard.com/report/github.com/relvacode/iso8601)


```
go get github.com/relvacode/iso8601
```

The built-in `RFC3339` time layout in Go is too restrictive to support ISO8601 date-times.

This library parses any ISO8601 date into a native Go time object without regular expressions.

## Usage

```go
import "github.com/relvacode/iso8601"

// iso8601 can be used as a drop-in replacement for time.Time with JSON responses
type ExternalAPIResponse struct {
	Timestamp *iso8601.Time
}


func main() {
	// ParseString can also be called directly
	t, err := iso8601.ParseString("2020-01-02T16:20:00")
}
```

## Benchmark

```
BenchmarkParse-16        	13364954	        77.7 ns/op	       0 B/op	       0 allocs/op
```

