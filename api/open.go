package api

import (
	"fmt"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) open(r rpc.Responder, c *rpc.Call) {
	var args []string
	c.Receive(&args)

	f, err := s.task.NS().Open(args[0])
	if err != nil {
		r.Return(err)
		return
	}

	fd := s.task.OpenFD(f, args[0])
	r.Return(uint64(fd))
}

func (s *syscaller) create(r rpc.Responder, c *rpc.Call) {
	var args []string
	c.Receive(&args)

	f, err := fs.Create(s.task.NS(), args[0])
	if err != nil {
		r.Return(err)
		return
	}

	fd := s.task.OpenFD(f, args[0])
	r.Return(uint64(fd))
}

func (s *syscaller) openFile(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 3); err != nil {
		r.Return(fmt.Errorf("openFile: %w", err))
		return
	}

	path, err := stringArg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("openFile: arg 0: %w", err))
		return
	}

	flags, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("openFile: arg 1: %w", err))
		return
	}

	mode, err := uint64Arg(args[2])
	if err != nil {
		r.Return(fmt.Errorf("openFile: arg 2: %w", err))
		return
	}

	f, err := fs.OpenFile(s.task.NS(), path, int(flags), fs.FileMode(mode))
	if err != nil {
		r.Return(err)
		return
	}

	fd := s.task.OpenFD(f, path)
	r.Return(uint64(fd))
}
