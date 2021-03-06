// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package data implements support for loading module data.
//
// The module data must be stored in the data directory, in the module root
// directory.
//
// The data may be provided by the main module or one of the active modules,
// and it will be accessed by the main module.
package data

import (
	"io"
	"os"
)

// Locator is responsible for finding how to load module data.
type Locator interface {
	// Locate returns the data loader for the module named by modpath.
	Locate(modpath string) (Loader, error)

	// Name returns the locator name.
	Name() string
}

// Loader is responsible for loading a module data file.
type Loader interface {
	// Load returns the file associated at path.
	//
	// path must be a relative path, without the "data/" prefix.
	Load(path string) (File, error)

	// Module returns the module the loader is associated with.
	Module() *Module
}

// File represents a module data file.
type File interface {
	// Name return the file name, relative to the data directory.
	Name() string

	// Path returns the absolute path to the file.
	//
	// If the file is embedded, Path will return an empty string.
	Path() string

	// Lstat returns information about the file.  If the file is a symbolic
	// link, Lstat returns information about the link itself, not the file it
	// points to.
	Lstat() (os.FileInfo, error)

	// Open provides access to the data within a regular file.  Open may return
	// an error if called on a directory or symbolic link.
	Open() (io.ReadCloser, error)
}

// DefaultLocator is the default locator.  It is automatically set to, in
// order:
//
//  1. If build info is not available, the "null" locator
//  2. If the main module version is "(devel)", the "fs:gopath" locator
//  3. The "fs:user" locator, if the main module is in $GODATA
//  4. The "fs:modcache" locator, if the main module is in the module cache
//  5. The "null" locator
var DefaultLocator Locator

// locators is a map with available locators.
var locators map[string]Locator

// Locate returns the loader for the main module, using the default locator.
func Locate() (Loader, error) {
	if info == nil {
		// The DefaultLocator is "null", if build info is not available.
		return DefaultLocator.Locate("")
	}

	return DefaultLocator.Locate(info.Main.Path)
}

// LocatorByName returns the locator by its name, or nil if not available.
//
// LocatorByName may return the "null" locator, if the requested locator is not
// usable on the host system, e.g. if the GOPATH environment variable is not
// defined for the "fs:gopath" locator.
//
// Supported locators are "fs:gopath", "fs:modcache" and "fs:user".
func LocatorByName(name string) Locator {
	l, ok := locators[name]
	if !ok {
		return nil
	}

	return l
}

// Load returns the file associated at path for the main module, using the
// default locator.
//
// path must be a relative path, without the "data/" prefix.
func Load(path string) (File, error) {
	l, err := Locate()
	if err != nil {
		return nil, err
	}

	return l.Load(path)
}

func init() {
	// Ensure the global info variable is initialized before everything else.
	bi, ok := readBuildInfo()
	if !ok {
		return
	}
	info = bi

	// Initialize the supported locators.
	locators = map[string]Locator{
		"fs:gopath":   newGopathLocator(),
		"fs:modcache": newModcacheLocator(),
		"fs:user":     newUserLocator(),
	}
}

func init() {
	DefaultLocator = defaultLocator()
}
