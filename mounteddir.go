// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package fsex

import (
	"os"
	"path/filepath"
	"time"

	"github.com/vedranvuk/fs" // transitional package.
)

// MountedDir implements a read-only FS in a mounted directory.
type MountedDir struct {
	// root is the path to the mounted directory.
	root string
}

// NewMountedDir returns a new MountedDir instance.
func NewMountedDir(root string) fs.ReadDirFS { return &MountedDir{root: root} }

// Open implements fs.FS.
func (md *MountedDir) Open(filename string) (fs.File, error) {
	var fn string
	if fn = filepath.Join(md.root, filename); fs.ValidPath(filename) {
		f, err := os.Open(fn)
		if err != nil {
			return nil, err
		}
		return &file{f, fn}, nil
	}
	return nil, &fs.PathError{Op: "open", Path: fn, Err: fs.ErrInvalid}
}

// ReadDir implements fs.ReadDirFS.
func (md *MountedDir) ReadDir(name string) ([]fs.DirEntry, error) {
	f, err := os.Open(md.root)
	if err != nil {
		return nil, err
	}
	fis, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	des := make([]fs.DirEntry, 0, len(fis))
	for _, fi := range fis {
		des = append(des, &fileInfo{fi})
	}
	return des, nil
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
func (f *file) Read(b []byte) (int, error) {
	return f.file.Read(b)
}

// Write could implement fs.File.Write.
func (f *file) Write(b []byte) (int, error) {
	return f.file.Write(b)
}

// Seek could implement fs.File.Seek.
func (f *file) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

// Close implements fs.File.Close.
func (f *file) Close() error {
	return f.file.Close()
}

// ReadDir implements fs.ReadDirFile.
func (f *file) ReadDir(n int) ([]fs.DirEntry, error) {
	fis, err := f.file.Readdir(n)
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