# iso8601
A fast ISO8601 date parser for Go

[![Build Status](https://travis-ci.org/relvacode/iso8601.svg?branch=master)](https://travis-ci.org/relvacode/iso8601) [![Go Report Card](https://goreportcard.com/badge/github.com/relvacode/iso8601)](https://goreportcard.com/report/github.com/relvacode/iso8601)

```go
import "github.com/relvacode/iso8601"
```

When working with dates in Go, especially with API communication the default `RFC3339` time layout is too restrictive to support the wide range of dates supported in the ISO8601 specification.

This library intends to parse any date that looks like the ISO8601 standard into native Go time.

## Performance

This library is efficient with no allocations needed to parse a full date.

    BenchmarkParse-8        20000000               100 ns/op               0 B/op          0 allocs/op
