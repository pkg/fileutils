package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	var testroot string // unique root, populated on each test run

	// joinpath turns relative paths into paths abosolute to the test root.
	joinpath := func(args ...string) string {
		return filepath.Join(append([]string{testroot}, args...)...)
	}

	// mkdir creates a directory inside testroot
	mkdir := func(perm os.FileMode, path ...string) {
		if err := os.Mkdir(joinpath(path...), perm); err != nil {
			t.Fatal(err)
		}
	}

	// mkfile creates a file with the specified contents inside testroot
	mkfile := func(perm os.FileMode, contents string, path ...string) {
		if err := ioutil.WriteFile(joinpath(path...), []byte(contents), perm); err != nil {
			t.Fatal(err)
		}
	}

	pass := func(t *testing.T, src, dst string, err error) {
		if err != nil {
			t.Errorf("CopyFile(%q, %q): got %v, expected %v", dst, src, err, nil)
		}
	}

	tests := []struct {
		setup    func(t *testing.T)
		dst, src string // automatically joined to testroot
		check    func(t *testing.T, src, dst string, err error)
	}{{
		setup: func(*testing.T) {
			mkdir(0755, "a")
			mkfile(0644, "file1", "a", "file1")
		},
		dst:   "a/file2",
		src:   "a/file1",
		check: pass,
	}, {
		setup: func(*testing.T) {
			mkdir(0755, "a")
			mkfile(0644, "file1", "a", "file1")
			mkfile(0644, "file2", "a", "file2")
		},
		dst:   "a/file2",
		src:   "a/file1",
		check: pass,
	}}

	// use a single tmpdir as the root of all tests to avoid spewing a million
	// tempdirs into $TMP during the test or on failure. Also, this means not
	// having to handle the cleanup of each
	root, err := ioutil.TempDir("", "testcopyfile")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Fatal(err)
		}
	}()

	for i, tt := range tests {
		testroot, err = ioutil.TempDir(root, fmt.Sprintf("test-%d", i))
		if err != nil {
			t.Fatal(err)
		}
		tt.setup(t)
		src := joinpath(filepath.FromSlash(tt.src))
		dst := joinpath(filepath.FromSlash(tt.dst))

		err := CopyFile(dst, src)
		tt.check(t, src, dst, err)
	}
}
