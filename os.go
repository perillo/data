// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The UserDataDir function has been adapted from os.UserConfigDir.
// Copyright 2009 The Go Authors. All rights reserved.

package data

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

// AppName returns the Go application name.  It is derived from the main
// package import path, using the last path segment.
//
// If build info is not available, it returns an empty string.
func AppName() string {
	// TODO(mperillo): Implement AppName to return a Java package like path.
	if info == nil {
		return ""
	}

	idx := strings.LastIndexByte(info.Path, '/')
	if idx < 0 {
		return info.Path
	}

	return info.Path[idx+1:]
}

//
// UserDataDir returns the default root directory to use for user-specific
// data.  Users should create their own application-specific subdirectory within
// this one and use that.
//
// On Unix systems, it returns $XDG_DATA_HOME as specified by
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
// if non-empty, else $HOME/.local/share.
// On Darwin, it returns $HOME/Library/Application Support.
// On Windows, it returns %LocalAppData%.
// On Plan 9, it returns $home/lib.
//
// If the location cannot be determined (for example, $HOME is not defined),
// then it will return an error.
func UserDataDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LocalAppData")
		if dir == "" {
			return "", errors.New("%LocalAppData% is not defined")
		}

	case "darwin":
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}
		dir += "/Library/Application Support"

	case "plan9":
		dir = os.Getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}
		dir += "/lib"

	default: // Unix
		dir = os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_DATA_HOME nor $HOME are defined")
			}
			dir += "/.local/share"
		}
	}

	return dir, nil
}
