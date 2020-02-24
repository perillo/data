// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import "runtime/debug"

// info stores the value returned by readBuildInfo.  It can be nil.
var info *buildInfo

// buildInfo represents the build information read from the running binary.
type buildInfo struct {
	Path string   // The main package path
	Main Module   // The main module information
	Deps []Module // Module dependencies
}

// readBuildInfo calls runtime/debug.ReadBuildInfo and convert the data to our
// internal buildInfo.
func readBuildInfo() (*buildInfo, bool) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, false
	}

	info := &buildInfo{
		Path: bi.Path,
		Main: fromDebug(&bi.Main), // bi.Main is not a pointer, unlike bi.Deps.
		Deps: make([]Module, len(bi.Deps)),
	}
	for i, m := range bi.Deps {
		info.Deps[i] = fromDebug(m)
	}

	return info, true
}

// Module represents a module.
type Module struct {
	Path    string  // module path
	Version string  // module version
	Sum     string  // checksum
	Replace *Module // replaced by this module
}

// fromDebug converts the module from debug.Module to Module.
func fromDebug(m *debug.Module) Module {
	mod := Module{
		Path:    m.Path,
		Version: m.Version,
		Sum:     m.Sum,
	}
	if mod.Replace != nil {
		// Replace is not recursive.
		mod.Replace = &Module{
			Path:    m.Replace.Path,
			Version: m.Replace.Version,
			Sum:     m.Replace.Sum,
		}
	}

	return mod
}

// String implements the Stringer interface.
func (m *Module) String() string {
	s := m.Path
	if m.Version != "" {
		s += "@" + m.Version
	}

	return s
}
