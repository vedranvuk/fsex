// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package fsex

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/vedranvuk/fs" // transitional package.
)

var (
	// ErrOpNotSupported is returned when an FS operation is not supported.
	ErrOpNotSupported = errors.New("op not supported")
)

// MountedDir implements a read-only FS in a mounted directory.
type MountedDir struct {
	root string
}

// NewMountedDir returns a new MountedDir instance.
func NewMountedDir(root string) (fs.FS, error) {
	return &MountedDir{root: root}, nil
}

// Open implements fs.FS.
func (md *MountedDir) Open(filename string) (fs.File, error) {
	var fn string
	if fn = filepath.Join(md.root, filename); fs.ValidPath(fn) {
		f, err := os.Open(fn)
		if err != nil {
			return nil, err
		}
		return &file{f, fn}, nil
	}
	return nil, &fs.PathError{Op: "open", Path: fn, Err: fs.ErrInvalid}
}

// file implements fs.File.
type file struct {
	file *os.File
	path string
}

// Stat implements fs.File.Stat.
func (f *file) Stat() (fs.FileInfo, error) {
	fi, err := f.file.Stat()
	if err != nil {
		return nil, err
	}
	return &fileInfo{fi}, nil
}

// Read implements fs.File.Read.
func (f *file) Read(b []byte) (int, error) { return f.file.Read(b) }

// Write could implement fs.File.Write.
func (f *file) Write(b []byte) (int, error) {
	return 0, &os.PathError{Op: "write", Path: f.path, Err: ErrOpNotSupported}
}

// Seek could implement fs.File.Seek.
func (f *file) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

// Close implements fs.File.Close.
func (f *file) Close() error { return f.file.Close() }

// ReadDir implements fs.ReadDirFile.
func (f *file) ReadDir(name string) ([]fs.DirEntry, error) {
	fis, err := f.file.Readdir(-1)
	if err != nil {
		return nil, err
	}
	des := make([]fs.DirEntry, 0, len(fis))
	for _, fi := range fis {
		des = append(des, &fileInfo{fi})
	}
	return des, nil
}

// fileInfo implements fs.DirEntry and fs.FileInfo.
type fileInfo struct{ fi os.FileInfo }

// Name implements fs.DirEntry and fs.FileInfo.
func (fi *fileInfo) Name() string { return fi.fi.Name() }

// IsDir implements fs.DirEntry and fs.FileInfo.
func (fi *fileInfo) IsDir() bool { return fi.fi.IsDir() }

// Type implements fs.DirEntry.
func (fi *fileInfo) Type() fs.FileMode { return fs.FileMode(fi.fi.Mode()) }

// Info implements fs.DirEntry.
func (fi *fileInfo) Info() (fs.FileInfo, error) {
	info, err := os.Stat(fi.fi.Name())
	if err != nil {
		return nil, err
	}
	return &fileInfo{info}, nil
}

// Size implements fs.FileInfo.
func (fi *fileInfo) Size() int64 { return fi.fi.Size() }

// Mode implements fs.FileInfo.
func (fi *fileInfo) Mode() fs.FileMode { return fs.FileMode(fi.fi.Mode()) }

// ModTime implements fs.FileInfo.
func (fi *fileInfo) ModTime() time.Time { return fi.fi.ModTime() }

// Sys implements fs.FileInfo.
func (fi *fileInfo) Sys() interface{} { return fi.fi.Sys() }
