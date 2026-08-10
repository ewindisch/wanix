package api

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// CBOR-decoded arguments do not always arrive with the Go type a handler
// expects: integral numbers decode as uint64/int64, tagged values wrap their
// content, and encoders disagree on float vs int. These helpers coerce like
// bytesArg does. A failed type assertion in an RPC handler panics and kills
// the whole runtime, so handlers validate through these and return errors.

// argCount guards against indexing past the decoded argument list.
func argCount(args []any, want int) error {
	if len(args) < want {
		return fmt.Errorf("want %d args, got %d", want, len(args))
	}
	return nil
}

// stringArg validates a string argument.
func stringArg(value any) (string, error) {
	if tag, ok := value.(cbor.Tag); ok {
		return stringArg(tag.Content)
	}
	s, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("want string, got %T", value)
	}
	return s, nil
}

// uint64Arg coerces a numeric argument, rejecting negative and fractional
// values.
func uint64Arg(value any) (uint64, error) {
	if tag, ok := value.(cbor.Tag); ok {
		return uint64Arg(tag.Content)
	}
	switch n := value.(type) {
	case uint64:
		return n, nil
	case int64:
		if n < 0 {
			return 0, fmt.Errorf("want unsigned integer, got %d", n)
		}
		return uint64(n), nil
	case int:
		if n < 0 {
			return 0, fmt.Errorf("want unsigned integer, got %d", n)
		}
		return uint64(n), nil
	case float64:
		if n < 0 || n != float64(uint64(n)) {
			return 0, fmt.Errorf("want unsigned integer, got %v", n)
		}
		return uint64(n), nil
	}
	return 0, fmt.Errorf("want unsigned integer, got %T", value)
}

// unixSeconds coerces a timestamp argument into fractional seconds.
// Integer-valued timestamps decode as uint64/int64 rather than float64.
func unixSeconds(value any) (float64, error) {
	if tag, ok := value.(cbor.Tag); ok {
		return unixSeconds(tag.Content)
	}
	switch n := value.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case int:
		return float64(n), nil
	}
	return 0, fmt.Errorf("want numeric timestamp, got %T", value)
}
