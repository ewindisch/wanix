package api

import (
	"fmt"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) truncate(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("truncate: %w", err))
		return
	}

	path, err := stringArg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("truncate: arg 0: %w", err))
		return
	}

	usize, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("truncate: arg 1: %w", err))
		return
	}
	size := int64(usize)

	err = fs.Truncate(s.task.NS(), path, size)
	if err != nil {
		r.Return(err)
		return
	}
}

func (s *syscaller) ftruncate(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("ftruncate: %w", err))
		return
	}

	ufd, err := uint64Arg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("ftruncate: arg 0: %w", err))
		return
	}
	fd := int(ufd)

	usize, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("ftruncate: arg 1: %w", err))
		return
	}
	size := int64(usize)

	_, path, err := s.task.FD(fd)
	if err != nil {
		r.Return(err)
		return
	}

	err = fs.Truncate(s.task.NS(), path, size)
	if err != nil {
		r.Return(err)
		return
	}
}
