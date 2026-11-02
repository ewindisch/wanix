package wanix

import (
	"net"
	"testing"

	"github.com/hugelgupf/p9/p9"
	"tractor.dev/wanix/fs"
	"tractor.dev/wanix/fs/p9kit"
)

func TestTaskSelfUsesCallerContextThroughRootMount(t *testing.T) {
	root, err := NewRoot()
	if err != nil {
		t.Fatal(err)
	}
	root.Register("test", autoDriver(func(*Task) error { return nil }))

	parent, err := root.fsys.Alloc("test", root)
	if err != nil {
		t.Fatal(err)
	}
	child, err := root.fsys.Alloc("test", parent)
	if err != nil {
		t.Fatal(err)
	}

	for _, task := range []*Task{root, parent, child} {
		got, err := fs.ReadFile(task.NS(), "#task/self/id")
		if err != nil {
			t.Fatalf("task %s: %v", task.ID(), err)
		}
		if want := task.ID() + "\n"; string(got) != want {
			t.Errorf("task %s: #task/self/id = %q, want %q", task.ID(), got, want)
		}
	}

	left, right := net.Pipe()
	server := p9.NewServer(p9kit.Attacher(child.NS()))
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handle(left, left) }()
	client, err := p9kit.ClientFS(right, "")
	if err != nil {
		t.Fatal(err)
	}
	got, err := fs.ReadFile(client, "#task/self/id")
	if err != nil {
		t.Fatal(err)
	}
	if want := child.ID() + "\n"; string(got) != want {
		t.Fatalf("9P #task/self/id = %q, want %q", got, want)
	}
	right.Close()
	left.Close()
	if err := <-serverDone; err != nil {
		t.Fatalf("9P server: %v", err)
	}
}
