// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package fsex

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/vedranvuk/fs"
)

func check(err error) bool {
	if err != nil {
		panic(err)
	}
	return true
}

func CreateTestFiles() map[string]bool {
	var files = map[string]bool{
		"rootdir1":                         true,
		"rootdir2":                         true,
		"rootdir1/subdir1":                 true,
		"rootdir1/subdir2":                 true,
		"rootdir1/subdir1/subdirfile1.ext": false,
		"rootdir1/subdir1/subdirfile2.ext": false,
		"rootfile1.ext":                    false,
		"rootfile2.ext":                    false,
	}
	var name string
	var isdir bool
	var err error
	for name, isdir = range files {
		if isdir {
			check(os.MkdirAll(path.Join("testdata", name), 0755))
			continue
		} else {
			check(os.MkdirAll(path.Join("testdata", path.Dir(name)), 0755))
		}
		if _, err = os.OpenFile(path.Join("testdata", name), os.O_CREATE, 0644); check(err) {
		}
	}
	return files
}

func RemoveTestData() {
	os.RemoveAll("testdata")
}

func printDirEntry(entry fs.DirEntry) {
	fmt.Printf(`?? DirEntry
  Name: '%s'
  IsDir: '%t'
  Type: '%v'
`,
		entry.Name(),
		entry.IsDir(),
		entry.Type(),
	)
}

func printFileInfo(info fs.FileInfo) {
	fmt.Printf(`?? FileInfo 
  Name: '%s'
  IsDir: '%t'
  Mode: '%v'
  Size: '%d'
  Sys: '%v'
  ModTime: '%v'
`,
		info.Name(),
		info.IsDir(),
		info.Mode(),
		info.Size(),
		info.Sys(),
		info.ModTime(),
	)
}

func TestMountedDir(t *testing.T) {
	var data = CreateTestFiles()
	defer RemoveTestData()

	var err error
	var md fs.FS
	if md, err = NewMountedDir("testdata"); err != nil {
		t.Fatal(err)
	}

	var name string
	var isdir bool
	var entries []fs.DirEntry
	var entry fs.DirEntry
	var file fs.File
	var info fs.FileInfo
	for name, isdir = range data {
		if isdir {
			fmt.Printf("{} Dir: '%s'\n", name)
			if entries, err = md.(fs.ReadDirFS).ReadDir(name); err != nil {
				t.Fatal(err)
			}
			for _, entry = range entries {
				printDirEntry(entry)
			}
			continue
		}
		fmt.Printf(">> File: '%s'\n", name)
		if file, err = md.Open(name); err != nil {
			t.Fatal(err)
		}
		if info, err = file.Stat(); err != nil {
			t.Fatal(err)
		}
		printFileInfo(info)
	}
}
