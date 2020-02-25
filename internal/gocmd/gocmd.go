// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocmd

// Getenv returns the named Go environment variable.
//
// If key does not exist, Getenv returns an empty string.
func Getenv(key string) (string, error) {
	stdout, err := Invoke("env", "GOPATH")
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}
