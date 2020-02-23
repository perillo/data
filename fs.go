// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fs.go source file implements the Loader and File interface for files on the
// local filesystem.

package data

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// fsLoader implements a Loader that loads module data from the filesystem.
type fsLoader struct {
	lc   Locator
	mod  *Module
	root string // absolute path to the module data directory.
}

// Load implements the Loader interface.
func (l *fsLoader) Load(path string) (File, error) {
	fl, err := l.load(path)
	if err != nil {
		lc := l.lc.Name()
		ld := l.mod

		return nil, fmt.Errorf("data: %s: %s: load: %v", lc, ld, err)
	}

	return fl, nil
}

func (l *fsLoader) load(path string) (File, error) {
	if filepath.IsAbs(path) {
		return nil, fmt.Errorf("path %s is not a relative path", path)
	}
	if !isDir(l.root) {
		return nil, fmt.Errorf("module %v does not have data", l.mod)
	}

	// It is responsibility of File to report an error if path does not exists.
	file := &fsFile{
		lc:   l.lc,
		ld:   l,
		root: l.root,
		path: path,
	}

	return file, nil
}

// Module implements the Loader interface.
func (l *fsLoader) Module() *Module {
	return l.mod
}

// fsFile represents a file on the local filesystem.
type fsFile struct {
	lc   Locator
	ld   Loader
	root string // absolute path to the module data directory
	path string // path to the data file, relative to the module data directory
}

// Name implements the File interface.
func (f *fsFile) Name() string {
	return f.path
}

// Path implements the File interface.
func (f *fsFile) Path() string {
	return filepath.Join(f.root, f.path)
}

// Lstat implements the File interface.
func (f *fsFile) Lstat() (os.FileInfo, error) {
	path := f.Path()

	fi, err := os.Lstat(path)
	if err != nil {
		lc := f.lc.Name()
		ld := f.ld.Module()

		return nil, fmt.Errorf("data: %s: %s: lstat %s: %v", lc, ld, f.path, err)
	}

	return fi, nil
}

// Open implements the File interface.
func (f *fsFile) Open() (io.ReadCloser, error) {
	path := f.Path()

	// TODO(mperillo): Check that the file is a regular file, in Open.
	rc, err := os.Open(path)
	if err != nil {
		lc := f.lc.Name()
		ld := f.ld.Module()

		return nil, fmt.Errorf("data: %s: %s: open %s: %v", lc, ld, f.path, err)
	}

	return rc, nil
}

// isDir returns true if path exists and it is a directory.
func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}
