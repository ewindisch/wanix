package api

import (
	"fmt"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) chmod(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("chmod: %w", err))
		return
	}

	path, err := stringArg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("chmod: arg 0: %w", err))
		return
	}

	umode, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("chmod: arg 1: %w", err))
		return
	}
	mode := fs.FileMode(umode)

	err = fs.Chmod(s.task.NS(), path, mode)
	if err != nil {
		r.Return(err)
		return
	}
}

func (s *syscaller) fchmod(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("fchmod: %w", err))
		return
	}

	ufd, err := uint64Arg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("fchmod: arg 0: %w", err))
		return
	}
	fd := int(ufd)

	umode, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("fchmod: arg 1: %w", err))
		return
	}
	mode := fs.FileMode(umode)

	_, path, err := s.task.FD(fd)
	if err != nil {
		r.Return(err)
		return
	}

	err = fs.Chmod(s.task.NS(), path, mode)
	if err != nil {
		r.Return(err)
		return
	}
}
