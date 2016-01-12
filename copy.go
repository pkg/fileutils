package fileutils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// A Copier copies files.
// The operation of Copier's public functions are controled by its
// public fields. If none are set, the Copier behaves accoriding to
// the zero value rules of each public field.
type Copier struct {
}

// CopyFile copies the contents of src to dst atomically.
func (c *Copier) CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "copyfile")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	const perm = 0644
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), dst); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}

// CopyFile is a convenience method that calls CopyFile on a Copier
// zero value.
func CopyFile(dst, src string) error {
	var c Copier
	return c.CopyFile(dst, src)
}
