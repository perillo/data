// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"fmt"
)

// defaultLocator returns the default locator as specified in the
// documentation.
func defaultLocator() Locator {
	// Check the build info to determine if this executable was installed with
	// go get.
	if info == nil {
		return &nullLocator{
			err: errors.New("build info is not available"),
		}
	}

	if info.Main.Version == "(devel)" {
		// Development mode, try to use the "fs:gopath" locator.
		return LocatorByName("fs:gopath")
	}

	// Installed mode.  Determine if the data is in the user data directory or
	// in the module cache.
	if l := LocatorByName("fs:user"); l.Name() == "fs:user" {
		return l
	}
	if l := LocatorByName("fs:modcache"); l.Name() == "fs:modcache" {
		return l
	}

	// Fallback to the "null" locator.
	return &nullLocator{
		err: errors.New("no locator is available"),
	}
}

// find finds the module named by modpath in the build info.  find assumes that
// info is not nil.
func find(modpath string) (*Module, error) {
	// TODO(mperillo): Use a module cache in find.
	if modpath == info.Main.Path {
		return &info.Main, nil
	}

	// TODO(mperillo): Decide what to do if there are multiple versions of the
	// same module.  Currently we return the first version found.
	for _, mod := range info.Deps {
		if mod.Path != modpath {
			continue
		}
		if mod.Replace != nil {
			return mod.Replace, nil
		}

		return &mod, nil
	}

	return nil, fmt.Errorf("module %s is not an active module", modpath)
}

// nullLocator is a Locator that always return an error.
type nullLocator struct {
	err error
}

// Locate implements the Locator interface.
func (l *nullLocator) Locate(modpath string) (Loader, error) {
	return nil, mkerr(l, l.err)
}

// Name implements the Locator interface.
func (l *nullLocator) Name() string {
	return "null"
}
