// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gopath.go source file implements the "fs:gopath" locator.

package data

import (
	"fmt"
	"path/filepath"

	"github.com/perillo/data/internal/gocmd"
)

// gopathLocator implements the "fs:gopath" locator that locates a module in
// $GOPATH.
type gopathLocator struct {
	pathList []string
}

// newGopathLocator returns a new "fs:gopath" locator, for modules in $GOPATH.
//
// newGopathLocator returns the "null" locator if $GOPATH is not available.
func newGopathLocator() Locator {
	gopath, err := gopath()
	if err != nil {
		return &nullLocator{
			err: err,
		}
	}

	l := &gopathLocator{
		pathList: filepath.SplitList(gopath),
	}

	// Check if the main module is in $GOPATH.
	if _, err := l.locate(info.Main.Path); err != nil {
		return &nullLocator{
			err: fmt.Errorf("main module %s is not in $GOPATH", &info.Main),
		}
	}

	return l
}

// Locate implements the Locator interface.
func (l *gopathLocator) Locate(modpath string) (Loader, error) {
	ld, err := l.locate(modpath)
	if err != nil {
		return nil, mkerr(l, err)
	}

	return ld, nil
}

func (l *gopathLocator) locate(modpath string) (Loader, error) {
	// Find module in build info.
	mod, err := find(modpath)
	if err != nil {
		return nil, err
	}

	// Search the module path in $GOPATH.
	for _, root := range l.pathList {
		dirpath := filepath.Join(root, "src", mod.Path)
		if isDir(dirpath) {
			// It is responsibility of Loader to report an error if the data
			// directory does not exists.
			ld := &fsLoader{
				lc:   l,
				mod:  mod,
				root: filepath.Join(dirpath, "data"),
			}

			return ld, nil
		}

	}

	return nil, fmt.Errorf("module %s is not in $GOPATH", modpath)
}

// Name implements the Locator interface.
func (l gopathLocator) Name() string {
	return "fs:gopath"
}

// gopath returns the value of the GOPATH environment variable.
func gopath() (string, error) {
	value, err := gocmd.Getenv("GOPATH")
	if err != nil {
		return "", fmt.Errorf("GOPATH is not available: %v", err)
	}

	// If there is no error, value should not be empty.  But each directory in
	// the list may not exist.
	return value, nil
}
