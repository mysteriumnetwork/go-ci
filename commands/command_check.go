/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-ci" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
