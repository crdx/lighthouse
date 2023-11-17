package util

import (
	"io/fs"
	"path/filepath"
)

// PrefixFS is an fs.FS wrapper that prefixes calls to Open with Prefix.
type PrefixFS struct {
	FS     fs.FS
	Prefix string
}

func (self *PrefixFS) Open(name string) (fs.File, error) {
	return self.FS.Open(filepath.Join(self.Prefix, name))
}
