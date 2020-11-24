// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package fsex

import "github.com/vedranvuk/fs"

// RWFile is a fs.File with Write and Seek methods.
type RWFile interface {
	fs.File
	Write([]byte) (int, error)
	Seek(int64, int) (int64, error)
}
