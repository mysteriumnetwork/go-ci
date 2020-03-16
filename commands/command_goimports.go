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
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-ci/shell"
	"github.com/mysteriumnetwork/go-ci/util"
)

// Fetches the goimports binary
func GetImports() error {
	path, _ := util.GetGoBinaryPath("goimports")
	if path != "" {
		fmt.Println("Tool 'goimports' already installed")
		return nil
	}
	err := sh.RunV("go", "get", "golang.org/x/tools/cmd/goimports")
	if err != nil {
		fmt.Println("Could not go get goimports")
		return err
	}
	return nil
}

// GoImports checks for issues with go imports
func GoImports(pathToCheck string, excludes ...string) error {
	mg.Deps(GetImports)

	goimportsBinaryPath, err := util.GetGoBinaryPath("goimports")
	if err != nil {
		fmt.Println("Tool 'goimports' not found")
		return err
	}
	gopath := util.GetGoPath()
	dirs, err := util.GetPackagePathsWithExcludes(pathToCheck, excludes...)
	if err != nil {
		fmt.Println("go list crashed")
		return err
	}

	dirsToLook := make([]string, 0)
	for _, dir := range dirs {
		absolutePath := path.Join(gopath, "src", dir)
		res, _ := ioutil.ReadDir(absolutePath)
		for _, v := range res {
			if v.IsDir() {
				continue
			}
			extension := filepath.Ext(v.Name())
			if extension != ".go" {
				continue
			}
			path := path.Join(absolutePath, v.Name())
			dirsToLook = append(dirsToLook, path)
		}
	}

	args := []string{"-e", "-l"}
	args = append(args, dirsToLook...)
	out, err := sh.Output(goimportsBinaryPath, args...)
	if err != nil {
		fmt.Println("Could not run goimports")
		return err
	}
	if len(out) != 0 {
		fmt.Println("The following files contain go import errors:")
		fmt.Println(out)
		return errors.New("not all imports follow the goimports format")
	}
	fmt.Println("Goimports is happy - all files are OK!")
	return nil
}

// GoImportsD checks for issues with go imports.
//
// Instead of packages, it operates on directories, thus it is compatible with gomodules outside GOPATH.
//
// Example:
//     commands.GoImportsD(".", "docs")
func GoImportsD(dir string, excludes ...string) error {
	mg.Deps(GetImports)
	goimportsBin, err := util.GetGoBinaryPath("goimports")
	if err != nil {
		fmt.Println("Tool 'goimports' not found")
		return err
	}
	var allExcludes []string
	allExcludes = append(allExcludes, excludes...)
	allExcludes = append(allExcludes, util.GoLintExcludes()...)
	dirs, err := util.GetProjectFileDirectories(allExcludes)
	if err != nil {
		return err
	}
	out, err := shell.NewCmd(goimportsBin + " -e -l -d " + strings.Join(dirs, " ")).Output()
	if err != nil {
		fmt.Println("goimports: error executing")
		return err
	}
	if len(out) != 0 {
		fmt.Println("goimports: the following files contain go import errors:")
		fmt.Println(out)
		return errors.New("goimports: not all imports follow the goimports format")
	}
	fmt.Println("goimports: all files are OK!")
	return nil
}
