package config

import (
	"io/fs"
	"os"
	"time"
)

type mapFS struct {
	files map[string][]byte
}

func newMapFS() *mapFS {
	return &mapFS{files: make(map[string][]byte)}
}

func (m *mapFS) Stat(name string) (fs.FileInfo, error) {
	if _, ok := m.files[name]; ok {
		return &mapFileInfo{name: name}, nil
	}

	return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrNotExist}
}

func (m *mapFS) Create(name string) (*os.File, error) {
	m.files[name] = nil

	return nil, nil
}

func (m *mapFS) MkdirAll(path string, _ fs.FileMode) error {
	m.files[path] = nil

	return nil
}

func (m *mapFS) ReadFile(name string) ([]byte, error) {
	content, ok := m.files[name]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}

	return content, nil
}

func (m *mapFS) writeFile(name string, content []byte) {
	m.files[name] = content
}

type mapFileInfo struct {
	name string
}

func (i *mapFileInfo) Name() string      { return i.name }
func (i *mapFileInfo) Size() int64       { return 0 }
func (i *mapFileInfo) Mode() fs.FileMode { return 0o444 }
func (i *mapFileInfo) ModTime() time.Time { return time.Time{} }
func (i *mapFileInfo) IsDir() bool       { return false }
func (i *mapFileInfo) Sys() any          { return nil }
