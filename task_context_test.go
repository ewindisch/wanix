package wanix

import (
	"testing"

	"tractor.dev/wanix/fs"
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
}
