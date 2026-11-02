package tarfs

import (
	"archive/tar"
	"bytes"
	"io/fs"
	"testing"
)

func TestFromPreservesHeadersAndReadsNestedFiles(t *testing.T) {
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	for _, h := range []*tar.Header{
		{Name: "usr/", Typeflag: tar.TypeDir, Mode: 0755},
		{Name: "usr/bin/", Typeflag: tar.TypeDir, Mode: 0755},
		{Name: "usr/bin/tool", Typeflag: tar.TypeReg, Mode: 0755, Size: 4},
	} {
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if h.Size > 0 {
			if _, err := tw.Write([]byte("wasm")); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := From(tar.NewReader(bytes.NewReader(archive.Bytes())))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"usr", "usr/bin"} {
		info, err := fs.Stat(got, name)
		if err != nil {
			t.Fatalf("Stat(%q): %v", name, err)
		}
		if !info.IsDir() {
			t.Errorf("Stat(%q) mode = %v, want directory", name, info.Mode())
		}
	}
	data, err := fs.ReadFile(got, "usr/bin/tool")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "wasm" {
		t.Fatalf("tool = %q, want wasm", data)
	}
}
