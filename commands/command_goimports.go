// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

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

// GetImports installs goimports binary.
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

// GoImports checks for issues with go imports.
func GoImports(pathToCheck string, excludes ...string) error {
	mg.Deps(GetImports)

	goimportsBinaryPath, err := util.GetGoBinaryPath("goimports")
	if err != nil {
		fmt.Println("❌ GoImports")
		fmt.Println("Tool 'goimports' not found")
		return err
	}
	gopath := util.GetGoPath()
	dirs, err := util.GetPackagePathsWithExcludes(pathToCheck, excludes...)
	if err != nil {
		fmt.Println("❌ GoImports")
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
		fmt.Println("❌ GoImports")
		fmt.Println("Could not run goimports")
		return err
	}
	if len(out) != 0 {
		fmt.Println("The following files contain go import errors:")
		fmt.Println(out)
		return errors.New("not all imports follow the goimports format")
	}
	fmt.Println("✅ GoImports")
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
		fmt.Println("❌ GoImports")
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
		fmt.Println("❌ GoImports")
		fmt.Println("goimports: error executing")
		return err
	}
	if len(out) != 0 {
		fmt.Println("❌ GoImports")
		fmt.Println("goimports: the following files contain go import errors:")
		fmt.Println(out)
		return errors.New("goimports: not all imports follow the goimports format")
	}
	fmt.Println("✅ GoImports")
	return nil
}
