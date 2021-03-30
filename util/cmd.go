// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package util

import (
	"go/build"
	"os"
	"path"

	"github.com/magefile/mage/sh"
)

// GetGoPath returns the go path
func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

// GetGoBinaryPath looks for the given binary in path, if not checks if it's in $GOPATH/bin
func GetGoBinaryPath(binaryName string) (string, error) {
	res, err := sh.Output("which", binaryName)
	if err == nil {
		return res, nil
	}
	gopath := GetGoPath()
	binaryUnderGopath := path.Join(gopath, "bin", binaryName)
	if _, err := os.Stat(binaryUnderGopath); os.IsNotExist(err) {
		return "", err
	}
	return binaryUnderGopath, nil
}
