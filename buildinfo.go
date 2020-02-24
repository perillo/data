// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import "runtime/debug"

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
