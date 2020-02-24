// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// modcache.go source file implements the "fs:modcache" locator.

package data

import (
	"errors"
	"fmt"
	"path/filepath"
)

// modcacheLocator implements the "fs:modcache" locator that locates a module
// in the Go module cache.
type modcacheLocator struct {
	path string
}

// newModcacheLocator returns a new "fs:modcache" locator, for modules in the
// Go module cache.
//
// newModcacheLocator returns the "null" locator if $GOPATH is not available or
// the main module is not stored in the Go module cache.
func newModcacheLocator() Locator {
	gocache, err := gocache()
	if err != nil {
		return &nullLocator{
			err: err,
		}
	}

	// Check if the main module is cached.
	if !isCached(gocache, &info.Main) {
		return &nullLocator{
			err: fmt.Errorf("main module %s is not in the cache", &info.Main),
		}
	}

	return &modcacheLocator{
		path: gocache,
	}
}

// Locate implements the Locator interface.
func (l *modcacheLocator) Locate(modpath string) (Loader, error) {
	fl, err := l.locate(modpath)
	if err != nil {
		return nil, mkerr(l, err)
	}

	return fl, nil
}

func (l *modcacheLocator) locate(modpath string) (Loader, error) {
	// Find module in build info.
	mod, err := find(modpath)
	if err != nil {
		return nil, err
	}

	// Check that the module is in the cache.
	dirpath := filepath.Join(l.path, mod.String())
	if !isDir(dirpath) {
		return nil, fmt.Errorf("module %s is not in the cache", mod)
	}

	// It is responsibility of Loader to report an error if the data directory
	// does not exists.
	lf := &fsLoader{
		lc:   l,
		mod:  mod,
		root: filepath.Join(dirpath, "data"),
	}

	return lf, nil
}

// Name implements the Locator interface.
func (l modcacheLocator) Name() string {
	return "fs:modcache"
}

// gocache returns the path to the module cache.
func gocache() (string, error) {
	gopath, err := gopath()
	if err != nil {
		return "", err
	}

	// The module cache is in the first entry of $GOPATH.
	root := filepath.SplitList(gopath)[0]
	path := filepath.Join(root, "pkg", "mod")
	if !isDir(path) {
		// This may improve the user experience, since it easy to clean the
		// module cache with go clean -modcache.
		return "", errors.New("module cache is not available")
	}

	return path, nil
}

// isCached returns true if the module, at the specified version, is available
// in the module cache at gocache.
func isCached(gocache string, mod *Module) bool {
	path := filepath.Join(gocache, mod.String())

	return isDir(path)
}
