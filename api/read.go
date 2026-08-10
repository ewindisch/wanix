package api

import (
	"fmt"
	"io"

	"tractor.dev/toolkit-go/duplex/rpc"
)

func (s *syscaller) read(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("read: %w", err))
		return
	}

	ufd, err := uint64Arg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("read: arg 0: %w", err))
		return
	}

	f, _, err := s.task.FD(int(ufd))
	if err != nil {
		r.Return(err)
		return
	}

	count, err := uint64Arg(args[1])
	if err != nil {
		r.Return(fmt.Errorf("read: arg 1: %w", err))
		return
	}

	buf := make([]byte, count)
	n, err := f.Read(buf)
	if err == io.EOF {
		r.Return(nil)
		return
	}
	if err != nil {
		r.Return(err)
		return
	}

	r.Return(buf[:n])
}
