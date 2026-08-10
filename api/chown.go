package api

import (
	"fmt"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) chown(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 3); err != nil {
		r.Return(fmt.Errorf("chown: %w", err))
		return
	}

	path, err := stringArg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("chown: arg 0: %w", err))
		return
	}

	uuid, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("chown: arg 1: %w", err))
		return
	}
	uid := int(uuid)

	ugid, err := uint64Arg(args[2])
	if err != nil {
		r.Return(fmt.Errorf("chown: arg 2: %w", err))
		return
	}
	gid := int(ugid)

	err = fs.Chown(s.task.NS(), path, uid, gid)
	if err != nil {
		r.Return(err)
		return
	}
}

func (s *syscaller) fchown(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 3); err != nil {
		r.Return(fmt.Errorf("fchown: %w", err))
		return
	}

	ufd, err := uint64Arg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("fchown: arg 0: %w", err))
		return
	}
	fd := int(ufd)

	uuid, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("fchown: arg 1: %w", err))
		return
	}
	uid := int(uuid)

	ugid, err := uint64Arg(args[2])
	if err != nil {
		r.Return(fmt.Errorf("fchown: arg 2: %w", err))
		return
	}
	gid := int(ugid)

	_, path, err := s.task.FD(fd)
	if err != nil {
		r.Return(err)
		return
	}

	err = fs.Chown(s.task.NS(), path, uid, gid)
	if err != nil {
		r.Return(err)
		return
	}
}
