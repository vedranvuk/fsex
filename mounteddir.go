// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package fsex

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/vedranvuk/fs" // transitional package.
)

var (
	// ErrNotADir is returned when passing a file to NewMountedDir.
	ErrNotADir = fmt.Errorf("%w: not a dir", ErrFSEX)
)

// MountedDir implements a read-only FS in a mounted directory.
type MountedDir struct {
	// root is the path to the mounted directory.
	root string
}

// NewMountedDir returns a new MountedDir instance as an fs.FS from specified
// root directory which must be a path to a directory.
// If an error occurs it is returned.
func NewMountedDir(root string) (fs.FS, error) { 
	var md = &MountedDir{}
	var err error
	if md.root, err = filepath.Abs(root); err != nil {
		return nil, err
	}
	var info os.FileInfo
	if info, err = os.Stat(root); err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrNotADir
	}
	return md, nil
}

// Open implements fs.FS.
// Open opens the underlying file which must exist and not be a directory for
// reading and writing with ModePerm. If an error occurs it is returned.
func (md *MountedDir) Open(filename string) (fs.File, error) {
	var fn string
	var f *os.File
	var err error
	if fn = filepath.Join(md.root, filename); fs.ValidPath(filename) {
		if f, err = os.OpenFile(fn, os.O_RDWR, os.ModePerm); err != nil {
			return nil, err
		}
		return &file{f, fn}, nil
	}
	return nil, &fs.PathError{Op: "open", Path: fn, Err: fs.ErrInvalid}
}

// ReadDir implements fs.ReadDirFS.
func (md *MountedDir) ReadDir(name string) ([]fs.DirEntry, error) {
	var f *os.File
	var err error
	if f, err = os.Open(md.root); err != nil {
		return nil, err
	}
	var infos []os.FileInfo
	if infos, err = f.Readdir(-1); err != nil {
		return nil, err
	}
	var entries = make([]fs.DirEntry, 0, len(infos))
	for _, fi := range infos {
		entries = append(entries, &fileInfo{fi, path.Join(md.root, fi.Name())})
	}
	return entries, nil
}

// Glob implmenets fs.GlobFS.
func (md *MountedDir) Glob(pattern string) (matches []string, err error) {
	// Matches are made against absolute paths.
	if matches, err = filepath.Glob(filepath.Join(md.root, pattern)); err != nil {
		return nil, err
	}
	// Absolute path prefix to matches need to be stripped.
	var index int
	var match string
	for index, match = range matches {
		matches[index] = strings.TrimPrefix(match, md.root)[1:]
	}
	return
}

// file implements fs.File.
type file struct {
	file *os.File
	name string
}

// Stat implements fs.File.Stat.
func (f *file) Stat() (fs.FileInfo, error) {
	var info os.FileInfo
	var err error
	if info, err = f.file.Stat(); err != nil {
		return nil, err
	}
	return &fileInfo{info, f.name}, nil
}

// Read implements fs.File.Read.
func (f *file) Read(b []byte) (int, error) { return f.file.Read(b) }

// Write could implement fs.File.Write.
func (f *file) Write(b []byte) (int, error) { return f.file.Write(b) }

// Seek could implement fs.File.Seek.
func (f *file) Seek(offset int64, whence int) (int64, error) { 
	return f.file.Seek(offset, whence) 
}

// Close implements fs.File.Close.
func (f *file) Close() error { return f.file.Close() }

// ReadDir implements fs.ReadDirFile.
func (f *file) ReadDir(n int) ([]fs.DirEntry, error) {
	var infos []os.FileInfo
	var err error
	if infos, err = f.file.Readdir(n); err != nil {
		return nil, err
	}
	var entries = make([]fs.DirEntry, 0, len(infos))
	var info os.FileInfo
	for _, info = range infos {
		entries = append(entries, &fileInfo{info, path.Join(f.name, info.Name())})
	}
	return entries, nil
}

// fileInfo implements fs.DirEntry and fs.FileInfo.
type fileInfo struct{ 
	info os.FileInfo 
	name string
}

// Name implements fs.DirEntry and fs.FileInfo.
func (fi *fileInfo) Name() string { return fi.info.Name() }

// IsDir implements fs.DirEntry and fs.FileInfo.
func (fi *fileInfo) IsDir() bool { return fi.info.IsDir() }

// Type implements fs.DirEntry.
func (fi *fileInfo) Type() fs.FileMode { return fs.FileMode(fi.info.Mode()) }

// Info implements fs.DirEntry.
func (fi *fileInfo) Info() (fs.FileInfo, error) {
	info, err := os.Stat(fi.info.Name())
	if err != nil {
		return nil, err
	}
	fi.info = info
	return &fileInfo{info, fi.name}, nil
}

// Size implements fs.FileInfo.
func (fi *fileInfo) Size() int64 { return fi.info.Size() }

// Mode implements fs.FileInfo.
func (fi *fileInfo) Mode() fs.FileMode { return fs.FileMode(fi.info.Mode()) }

// ModTime implements fs.FileInfo.
func (fi *fileInfo) ModTime() time.Time { return fi.info.ModTime() }

// Sys implements fs.FileInfo.
func (fi *fileInfo) Sys() interface{} { return fi.info.Sys() }
