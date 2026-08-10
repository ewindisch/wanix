package api

import (
	"fmt"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) appendFile(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	// log.Println("WriteFile", args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("appendFile: %w", err))
		return
	}

	name, err := stringArg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("appendFile: arg 0: %w", err))
		return
	}

	data, err := bytesArg(args[1])
	if err != nil {
		r.Return(err)
		return
	}

	err = fs.AppendFile(s.task.NS(), name, data)
	if err != nil {
		r.Return(err)
		return
	}
}
