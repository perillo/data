// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gocmd implements a simple wrapper for cmd/go.
package gocmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Error is the error returned when the go command returns an error.
type Error struct {
	Argv   []string // arguments to the go command
	Stderr []byte   // the entire content of the go command stderr
	Err    error    // the original error from os/exec.Command.Run
}

// Error implements the error interface.
func (e *Error) Error() string {
	argv := strings.Trim(fmt.Sprint(e.Argv), "[]")
	stderr := string(e.Stderr)
	msg := "go " + argv + ": " + e.Err.Error()

	if stderr == "" {
		return msg
	}

	return msg + ": " + stderr
}

// Unwrap implements the Wrapper interface.
func (e *Error) Unwrap() error {
	return e.Err
}

// Invoke invokes a go command and return the stdout content, with whitespace
// trimmed.
//
// In case the go command exits with a non 0 exit status, the error will
// contain the entire content of the go command stderr, with whitespace
// trimmed.
func Invoke(verb string, args ...string) ([]byte, error) {
	args = append([]string{verb}, args...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command("go", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		err := &Error{
			Argv:   args,
			Stderr: normalize(stderr),
			Err:    err,
		}

		return nil, err
	}

	return normalize(stdout), nil
}

// normalize returns the data buffered in b with leading and trailing white
// space removed.
func normalize(b *bytes.Buffer) []byte {
	return bytes.TrimSpace(b.Bytes())
}
