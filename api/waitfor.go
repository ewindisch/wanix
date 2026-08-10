package api

import (
	"fmt"
	"time"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) waitFor(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("waitFor: %w", err))
		return
	}

	path, err := stringArg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("waitFor: arg 0: %w", err))
		return
	}

	timeout, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("waitFor: arg 1: %w", err))
		return
	}

	deadline := time.Now().Add(time.Duration(timeout) * time.Millisecond)
	for time.Now().Before(deadline) {
		exists, err := fs.Exists(s.task.NS(), path)
		if err != nil {
			r.Return(err)
			return
		}
		if exists {
			r.Return()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	r.Return(fs.ErrNotExist)
}
