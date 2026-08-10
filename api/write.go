package api

import (
	"fmt"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) write(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("write: %w", err))
		return
	}

	ufd, err := uint64Arg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("write: arg 0: %w", err))
		return
	}

	f, _, err := s.task.FD(int(ufd))
	if err != nil {
		r.Return(err)
		return
	}

	data, err := bytesArg(args[1])
	if err != nil {
		r.Return(err)
		return
	}

	n, err := fs.Write(f, data)
	if err != nil {
		r.Return(err)
		return
	}

	r.Return(uint64(n))
}

func (s *syscaller) writeAt(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 3); err != nil {
		r.Return(fmt.Errorf("writeAt: %w", err))
		return
	}

	ufd, err := uint64Arg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("writeAt: arg 0: %w", err))
		return
	}

	f, _, err := s.task.FD(int(ufd))
	if err != nil {
		r.Return(err)
		return
	}

	data, err := bytesArg(args[1])
	if err != nil {
		r.Return(err)
		return
	}

	offset, err := uint64Arg(args[2])
	if err != nil {
		r.Return(fmt.Errorf("writeAt: arg 2: %w", err))
		return
	}

	n, err := fs.WriteAt(f, data, int64(offset))
	if err != nil {
		r.Return(err)
		return
	}

	r.Return(n)
}
