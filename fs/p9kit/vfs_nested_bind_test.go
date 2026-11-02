package p9kit

import (
	"context"
	"testing"
	"testing/fstest"

	"tractor.dev/wanix/fs"
	"tractor.dev/wanix/fs/vfs"
)

func TestClientWalkNestedBindWithSynthesizedParents(t *testing.T) {
	root := fstest.MapFS{}
	payload := fstest.MapFS{"ferrule-rc.wasm": {Data: []byte("wasm payload")}}
	ns := vfs.New(context.Background())
	ns.Bind(root, ".", ".", fs.BindReplace)
	ns.Bind(payload, "ferrule-rc.wasm", "pkg/ferruleos-wasi-preview1/0.1.0/usr/bin/ferrule-rc.wasm", fs.BindReplace)

	client, cleanup := testSetup(t, ns)
	defer cleanup()
	data, err := fs.ReadFile(client, "pkg/ferruleos-wasi-preview1/0.1.0/usr/bin/ferrule-rc.wasm")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "wasm payload" {
		t.Fatalf("payload = %q", data)
	}
}
