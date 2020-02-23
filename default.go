// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// Module represents a module.
type Module struct {
	Path    string  // module path
	Version string  // module version
	Sum     string  // checksum
	Replace *Module // replaced by this module
}

// fromDebug converts the module from debug.Module to Module.
func fromDebug(m *debug.Module) *Module {
	mod := &Module{
		Path:    m.Path,
		Version: m.Version,
		Sum:     m.Sum,
	}
	if mod.Replace != nil {
		mod.Replace = fromDebug(m.Replace)
	}

	return mod
}

// String implements the Stringer interface.
func (m *Module) String() string {
	s := m.Path
	if m.Version != "" {
		s += " " + m.Version
	}
	if m.Replace != nil {
		s += " => " + m.Replace.Path
		if m.Replace.Version != "" {
			s += " " + m.Replace.Version
		}
	}

	return s
}

// buildInfo stores the value returned by ReadBuildInfo.  It can be nil.
var buildInfo *debug.BuildInfo

// defaultLocator returns the default locator as specified in the
// documentation.
func defaultLocator() Locator {
	// Read the build info to determine if this executable was installed with
	// go get.
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return &nullLocator{
			err: errors.New("build info is not available"),
		}
	}
	buildInfo = bi // cache the build info for later use

	if bi.Main.Version == "(devel)" {
		// Development mode, use the "fs:gopath" locator.  The implementation
		// assumes that Main.Path is in $GOPATH.
		return newGopathLocator()
	}

	// Installed mode.  Determine if the data is in the user data directory or
	// in the module cache.
	// TODO(mperillo): Implement the "fs:user" and "fs:modcache" locators.
	return &nullLocator{
		err: errors.New("only the \"fs:gopath\" locator is implemented"),
	}
}

// find finds the module named by modpath in the build info.  find assumes that
// buildInfo is not null.
func find(modpath string) (*Module, error) {
	// TODO(mperillo): Use a module cache in find.
	if modpath == buildInfo.Main.Path {
		// Note that BuildInfo.Main is not a pointer, unlike BuildInfo.Deps.
		return fromDebug(&buildInfo.Main), nil
	}

	// TODO(mperillo): Decide what to do if there are multiple versions of the
	// same module.  Currently we return the first version found.
	for _, mod := range buildInfo.Deps {
		if mod.Path != modpath {
			continue
		}
		if mod.Replace != nil {
			return fromDebug(mod.Replace), nil
		}

		return fromDebug(mod), nil
	}

	return nil, fmt.Errorf("module %s is not an active module", modpath)
}

// nullLocator is a Locator that always return an error.
type nullLocator struct {
	err error
}

// Locate implements the Locator interface.
func (l *nullLocator) Locate(modpath string) (Loader, error) {
	return nil, fmt.Errorf("data: %s: locate: %v", l.Name(), l.err)
}

// Name implements the Locator interface.
func (l *nullLocator) Name() string {
	return "null"
}
