// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// errors.go source file implements the Error type.

package data

import "fmt"

// Error records an error during a data operation.
//
// Error is returned by Locator.Locate, Loader.Load, File.Lstat and File.Open.
type Error struct {
	Locator Locator
	Loader  Loader
	File    File
	Op      string // the file operation used ("lstat" or "open")
	Err     error  // the underlying error
}

// Error implements the error interface.
func (e *Error) Error() string {
	msg := "data: " + e.Locator.Name() + ": "
	if e.Loader == nil {
		return msg + e.Err.Error()
	}

	msg = msg + e.Loader.Module().String() + ": "
	if e.File == nil {
		return msg + e.Err.Error()
	}

	return msg + e.Op + " " + e.File.Name() + ": " + e.Err.Error()
}

// Unwrap implements the Wrapper interface.
func (e *Error) Unwrap() error {
	return e.Err
}

// mkerr builds an error value from its arguments.  It will panic if there are
// no arguments.
//
// The mkerr function has been inspired by upspin.io/errors.  See
// https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html.
func mkerr(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to mkerr with no arguments")
	}

	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Locator:
			e.Locator = arg
		case Loader:
			e.Loader = arg
		case File:
			e.File = arg
		case string:
			e.Op = arg
		case error:
			e.Err = arg
		default:
			panicf("call to mkerr with unknown type %T, value %v", arg, arg)
		}
	}

	return e
}

func panicf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}
