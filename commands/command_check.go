// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package commands

import (
	"github.com/magefile/mage/mg"
)

// Check performs all the common checks
func Check(dir string, excludes ...string) error {
	copyrightWrapper := func() error {
		return Copyright(dir, excludes...)
	}
	importsWrapper := func() error {
		return GoImports(dir, excludes...)
	}
	goLintWrapper := func() error {
		return GoLint(dir, excludes...)
	}
	goVetWrapper := func() error {
		return GoVet(dir)
	}
	mg.Deps(copyrightWrapper, importsWrapper, goLintWrapper, goVetWrapper)
	return nil
}

// CheckD performs all the common checks on directories
// Instead of packages, it operates on directories, thus it is compatible with gomodules outside GOPATH.
func CheckD(dir string, excludes ...string) error {
	copyrightWrapper := func() error {
		return CopyrightD(dir, excludes...)
	}
	importsWrapper := func() error {
		return GoImportsD(dir, excludes...)
	}
	goLintWrapper := func() error {
		return GoLintD(dir, excludes...)
	}
	goVetWrapper := func() error {
		return GoVet(dir)
	}
	mg.Deps(copyrightWrapper, importsWrapper, goLintWrapper, goVetWrapper)
	return nil
}
