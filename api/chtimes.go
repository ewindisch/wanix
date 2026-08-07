package api

import (
	"fmt"
	"time"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

// unixSeconds coerces a CBOR-decoded timestamp into fractional seconds.
// Integer-valued timestamps decode as int64/uint64 rather than float64, and
// a failed type assertion here previously panicked and killed the runtime.
func unixSeconds(v any) (float64, error) {
	switch t := v.(type) {
	case float64:
		return t, nil
	case float32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case uint64:
		return float64(t), nil
	case int:
		return float64(t), nil
	}
	return 0, fmt.Errorf("expected numeric timestamp, got %T", v)
}

func (s *syscaller) chtimes(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	path, ok := args[0].(string)
	if !ok {
		r.Return(fmt.Errorf("chtimes: arg 0 is not a string"))
		return
	}

	// atime and mtime are in seconds (with fractional parts)
	atimeSec, err := unixSeconds(args[1])
	if err != nil {
		r.Return(fmt.Errorf("chtimes: arg 1: %w", err))
		return
	}
	atime := time.Unix(int64(atimeSec), int64((atimeSec-float64(int64(atimeSec)))*1e9))

	mtimeSec, err := unixSeconds(args[2])
	if err != nil {
		r.Return(fmt.Errorf("chtimes: arg 2: %w", err))
		return
	}
	mtime := time.Unix(int64(mtimeSec), int64((mtimeSec-float64(int64(mtimeSec)))*1e9))

	err = fs.Chtimes(s.task.NS(), path, atime, mtime)
	if err != nil {
		r.Return(err)
		return
	}
}
