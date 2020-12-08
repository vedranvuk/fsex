// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package fsex

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/vedranvuk/fs"
)

func cleanup() {
	os.Chdir("..")
	os.RemoveAll("testdata")
}

func createDir(filename string) {
	if err := os.MkdirAll(filename, 0755); err != nil {
		cleanup()
		panic(err)
	}
}
func createFile(filename string) {
	createDir(filepath.Dir(filename))
	file, err := os.OpenFile(filename, os.O_CREATE, 0644)
	if err != nil {
		cleanup()
		panic(err)
	}
	file.Close()
}

func createTestData() {
	createDir("testdata")
	os.Chdir("testdata")
	createFile("index.html")
	createFile("static/css/default.css")
	createFile("static/js/default.js")
}

func TestMountedDir(t *testing.T) {
	createTestData()
	defer cleanup()
	var fsys fs.FS
	var err error
	if fsys, err = NewMountedDir("."); err != nil {
		t.Fatal(err)
	}
	var file fs.File
	if file, err = fsys.Open("index.html"); err != nil {
		t.Fatal(err)
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}
	if file, err = fsys.Open("static/css/default.css"); err != nil {
		t.Fatal(err)
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}
	if file, err = fsys.Open("static/js/default.js"); err != nil {
		t.Fatal(err)
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}
	var infos []fs.DirEntry
	if infos, err = fsys.(fs.ReadDirFS).ReadDir("."); err != nil {
		t.Fatal(err)
	}
	var expected = map[string]bool{
		"static":     true,
		"index.html": true,
	}
	var ok bool
	for i := 0; i < len(infos); i++ {
		if _, ok = expected[infos[i].Name()]; !ok {
			t.Fatal("ReadDir failed.")
		}
		if testing.Verbose() {
			fmt.Printf("Name: '%s', IsDir: '%t'\n", infos[i].Name(), infos[i].IsDir())
		}
	}
	var matches []string
	if matches, err = fsys.(fs.GlobFS).Glob("*/*/*.js"); err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 || matches[0] != "static/js/default.js" {
		t.Fatal("Glob failed.")
	}
	if testing.Verbose() {
		for _, match := range matches {
			fmt.Printf("Match: '%s'\n", match)
		}
	}
}
