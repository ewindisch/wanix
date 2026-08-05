//go:build !js && !wasm

package headless_test

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"tractor.dev/wanix/fs"
	"tractor.dev/wanix/fs/fskit"
	"tractor.dev/wanix/fs/memfs"
	"tractor.dev/wanix/fs/pipe"
	"tractor.dev/wanix/headless"
)

// TestNativeGoJSRC is the real acceptance seam. It builds the in-tree rc
// module unless WANIX_RC_WASM selects an already-pinned release artifact.
func TestNativeGoJSRC(t *testing.T) {
	rcPath := os.Getenv("WANIX_RC_WASM")
	if rcPath == "" {
		rcPath = filepath.Join(t.TempDir(), "rc.wasm")
		build := exec.Command("go", "build", "-o", rcPath, "./cmd/rc")
		build.Dir = filepath.Join("..", "rc")
		build.Env = append(os.Environ(), "GOWORK=off", "GOOS=js", "GOARCH=wasm")
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build rc.wasm: %v\n%s", err, output)
		}
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("native GoJS acceptance requires node: %v", err)
	}
	rc, err := os.ReadFile(rcPath)
	if err != nil {
		t.Fatal(err)
	}
	root, err := headless.Bootstrap()
	if err != nil {
		t.Fatal(err)
	}
	bin := memfs.New()
	if err := root.NS().Bind(bin, ".", "bin"); err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile(root.NS(), "bin/rc.wasm", rc, 0755); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	sock := filepath.Join(dir, "api.sock")
	runner := filepath.Join("..", "native", "wanix-gojs-runner.mjs")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := headless.New(headless.Options{Root: root, APIAddr: "unix://" + sock, GoJSRunner: runner})
	done := make(chan error, 1)
	go func() { done <- srv.Run(ctx) }()
	for i := 0; i < 100 && !fileExists(sock); i++ {
		time.Sleep(time.Millisecond)
	}
	if !fileExists(sock) {
		t.Fatal("API socket was not created")
	}

	allocated, err := fs.OpenContext(context.Background(), root.NS(), "#task/new/auto")
	if err != nil {
		t.Fatal(err)
	}
	idBytes, err := io.ReadAll(allocated)
	_ = allocated.Close()
	if err != nil {
		t.Fatal(err)
	}
	id := strings.TrimSpace(string(idBytes))
	if id == "" {
		t.Fatal("task allocation returned no id")
	}
	stdin, _ := pipe.New(false)
	stdout, output := pipe.New(false)
	stderr, errors := pipe.New(false)
	fd := fskit.MapFS{
		"0": fskit.FileFS(&pipe.PortFile{Port: stdin, Name: "0"}, "0"),
		"1": fskit.FileFS(&pipe.PortFile{Port: stdout, Name: "1"}, "1"),
		"2": fskit.FileFS(&pipe.PortFile{Port: stderr, Name: "2"}, "2"),
	}
	if err := root.NS().Bind(fd, ".", "#task/"+id+"/fd"); err != nil {
		t.Fatal(err)
	}
	child, err := root.Lookup(id)
	if err != nil {
		t.Fatal(err)
	}
	if err := child.NS().Bind(fd, ".", "#task/"+id+"/fd"); err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile(root.NS(), "#task/"+id+"/cmd", []byte("bin/rc.wasm -c 'echo native'"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile(root.NS(), "#task/"+id+"/ctl", []byte("start"), 0644); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		exit, readErr := fs.ReadFile(root.NS(), "#task/"+id+"/exit")
		if readErr == nil && strings.TrimSpace(string(exit)) != "" {
			if strings.TrimSpace(string(exit)) != "0" {
				t.Fatalf("rc exit %q; stdout %q; stderr %q", exit, readPipe(output), readPipe(errors))
			}
			if !strings.Contains(readPipe(output), "native") {
				t.Fatalf("rc output %q, want native", readPipe(output))
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("rc did not exit; stderr %q", readPipe(errors))
}

// TestStartCommandWritesNativeStdout exercises the exact launch path used by
// `wanix headless --shell-command`. The other acceptance test binds an
// in-memory terminal explicitly; this one verifies that the native process
// stdout file is exposed to GoJS through Wanix's task FD table.
func TestStartCommandWritesNativeStdout(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("native GoJS acceptance requires node: %v", err)
	}
	rcPath := filepath.Join(t.TempDir(), "rc.wasm")
	build := exec.Command("go", "build", "-o", rcPath, "./cmd/rc")
	build.Dir = filepath.Join("..", "rc")
	build.Env = append(os.Environ(), "GOWORK=off", "GOOS=js", "GOARCH=wasm")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build rc.wasm: %v\n%s", err, output)
	}
	rc, err := os.ReadFile(rcPath)
	if err != nil {
		t.Fatal(err)
	}
	root, err := headless.Bootstrap()
	if err != nil {
		t.Fatal(err)
	}
	if err := root.NS().Bind(memfs.New(), ".", "bin"); err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile(root.NS(), "bin/rc.wasm", rc, 0755); err != nil {
		t.Fatal(err)
	}
	runner := filepath.Join("..", "native", "wanix-gojs-runner.mjs")
	srv := headless.New(headless.Options{Root: root, GoJSRunner: runner})
	if _, err := srv.RootTask(); err != nil {
		t.Fatal(err)
	}

	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = write
	t.Cleanup(func() { os.Stdout = original })
	task, err := headless.StartCommand(root, "bin/rc.wasm -c 'echo native-stdio'")
	if err != nil {
		t.Fatal(err)
	}
	code, err := headless.WaitExit(task, 20*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("rc exited %d", code)
	}
	if err := write.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
		t.Fatal(err)
	}
	output, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "native-stdio") {
		t.Fatalf("native stdout %q, want native-stdio", output)
	}
}

func fileExists(path string) bool { _, err := os.Stat(path); return err == nil }

func readPipe(port *pipe.Port) string {
	buf := make([]byte, port.Size())
	_, _ = port.Read(buf)
	return string(buf)
}
