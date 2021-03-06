package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Meat-Hook/migrate/core"
)

var _ core.FS = &FS{}

// FS implement core.FS.
type FS struct{}

// New returns new instance default filesystem.
func New() *FS {
	return &FS{}
}

// Open for implements core.FS.
func (F *FS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// Walk for implements core.FS.
func (F *FS) Walk(path string, cb func(string, fs.FileInfo) error) error {
	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("path [%s]: walk error: %w", path, err)
		}

		return cb(path, info)
	})
}

// Mkdir for implements core.FS.
func (F *FS) Mkdir(path string) error {
	return os.MkdirAll(path, 0700)
}

// SaveFile for implements core.FS.
func (F *FS) SaveFile(path string, buf []byte) error {
	return os.WriteFile(path, buf, fs.ModePerm)
}
