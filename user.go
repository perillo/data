// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// user.go source file implements the "fs:user" locator.

package data

import (
	"fmt"
	"os"
	"path/filepath"
)

// userLocator implements the "fs:user" locator that locates a module in the
// user data directory.
type userLocator struct {
	path string
}

// newUserLocator returns a new "fs:user" locator, for modules in the user data
// directory.
//
// newUserLocator returns the "null" locator if the user data directory is not
// available or the application is not stored in the user data directory.
func newUserLocator() Locator {
	godata, err := godata()
	if err != nil {
		return &nullLocator{
			err: err,
		}
	}

	l := &userLocator{
		path: godata,
	}

	// Check if the main module is in the user data directory.
	if _, err := l.locate(info.Main.Path); err != nil {
		return &nullLocator{
			err: fmt.Errorf("main module %s is not in the user data directory",
				&info.Main),
		}
	}

	return l
}

// Locate implements the Locator interface.
func (l *userLocator) Locate(modpath string) (Loader, error) {
	ld, err := l.locate(modpath)
	if err != nil {
		return nil, mkerr(l, err)
	}

	return ld, nil
}

func (l *userLocator) locate(modpath string) (Loader, error) {
	// Find module in build info.
	mod, err := find(modpath)
	if err != nil {
		return nil, err
	}

	var dirpath string
	if modpath == info.Main.Path {
		// The main module is special, and the data is stored in
		// $GODATA/$APPNAME.
		dirpath = filepath.Join(l.path, AppName())
	} else {
		// Active modules are stored in $GODATA/go-data, with the fully
		// versioned path flattened.
		dirpath = filepath.Join(l.path, "go-data", mod.FlatPath())
	}

	if !isDir(dirpath) {
		return nil, fmt.Errorf("module %s is not in user data directory", modpath)
	}

	// It is responsibility of Loader to report an error if the data
	// directory does not exists.
	ld := &fsLoader{
		lc:   l,
		mod:  mod,
		root: filepath.Join(dirpath, "data"),
	}

	return ld, nil
}

// Name implements the Locator interface.
func (l userLocator) Name() string {
	return "fs:user"
}

func godata() (string, error) {
	dir := os.Getenv("GODATA")
	if dir != "" {
		return dir, nil
	}

	return UserDataDir()
}
