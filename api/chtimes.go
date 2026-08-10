package api

import (
	"fmt"
	"time"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) chtimes(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 3); err != nil {
		r.Return(fmt.Errorf("chtimes: %w", err))
		return
	}

	path, err := stringArg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("chtimes: arg 0: %w", err))
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
