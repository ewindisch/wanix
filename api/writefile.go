package api

import (
	"fmt"

	"tractor.dev/toolkit-go/duplex/rpc"
	"tractor.dev/wanix/fs"
)

func (s *syscaller) writeFile(r rpc.Responder, c *rpc.Call) {
	var args []any
	c.Receive(&args)

	// log.Println("WriteFile", args)

	if err := argCount(args, 2); err != nil {
		r.Return(fmt.Errorf("writeFile: %w", err))
		return
	}

	name, err := stringArg(args[0])
	if err != nil {
		r.Return(fmt.Errorf("writeFile: arg 0: %w", err))
		return
	}

	data, err := bytesArg(args[1])
	if err != nil {
		r.Return(err)
		return
	}

	err = fs.WriteFile(s.task.NS(), name, data, 0x644)
	if err != nil {
		r.Return(err)
		return
	}
}
