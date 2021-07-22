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
	"regexp"

	"github.com/mysteriumnetwork/go-ci/util"
)

var copyrightRegex = regexp.MustCompile(`Copyright \([cC]\) \d{4}`)

func getFilesWithoutCopyright(dirsToCheck []string) ([]string, error) {
	badFiles := make([]string, 0)
	gopath := util.GetGoPath()

	for i := range dirsToCheck {
		absolutePath := path.Join(gopath, "src", dirsToCheck[i])
		files, err := ioutil.ReadDir(absolutePath)
		if err != nil {
			return badFiles, err
		}
		for j := range files {
			if files[j].IsDir() {
				continue
			}
			extension := filepath.Ext(files[j].Name())
			if extension != ".go" {
				continue
			}
			contents, err := ioutil.ReadFile(path.Join(absolutePath, files[j].Name()))
			if err != nil {
				return nil, err
			}
			match := copyrightRegex.Match(contents)
			if !match {
				badFiles = append(badFiles, path.Join(dirsToCheck[i], files[j].Name()))
			}
		}
	}
	return badFiles, nil
}

// Copyright checks for copyright headers in files
func Copyright(path string, excludes ...string) error {
	res, err := util.GetPackagePathsWithExcludes(path, excludes...)
	if err != nil {
		fmt.Println("❌ Copyright")
		fmt.Println("go list crashed")
		return err
	}
	badFiles, err := getFilesWithoutCopyright(res)
	if err != nil {
		fmt.Println("❌ Copyright")
		return err
	}
	if len(badFiles) != 0 {
		fmt.Println("❌ Copyright")
		fmt.Println("Following files are missing copyright headers:")
		for _, v := range badFiles {
			fmt.Println(v)
		}
		return errors.New("Missing copyright headers")
	}
	fmt.Println("✅ Copyright")
	return nil
}

// CopyrightD checks for copyright headers in files
//
// Instead of packages, it operates on directories, thus it is compatible with gomodules outside GOPATH.
//
// Example:
//     commands.CopyrightD(".", "docs")
func CopyrightD(path string, excludes ...string) error {
	var allExcludes []string
	allExcludes = append(allExcludes, excludes...)
	allExcludes = append(allExcludes, util.GoLintExcludes()...)
	res, err := util.GetProjectFileDirectories(allExcludes)
	if err != nil {
		fmt.Println("❌ Copyright")
		fmt.Println("copyright: go list crashed")
		return err
	}
	badFiles, err := getFilesWithoutCopyrightD(res)
	if err != nil {
		fmt.Println("❌ Copyright")
		fmt.Println("copyright: error listing files")
		return err
	}
	if len(badFiles) != 0 {
		fmt.Println("❌ Copyright")
		fmt.Println("copyright: following files are missing copyright headers:")
		for _, v := range badFiles {
			fmt.Println(v)
		}
		return errors.New("copyright: missing copyright headers")
	}
	fmt.Println("✅ Copyright")
	return nil
}

func getFilesWithoutCopyrightD(dirsToCheck []string) ([]string, error) {
	badFiles := make([]string, 0)

	for i := range dirsToCheck {
		files, err := ioutil.ReadDir(dirsToCheck[i])
		if err != nil {
			return badFiles, err
		}
		for j := range files {
			if files[j].IsDir() {
				continue
			}
			extension := filepath.Ext(files[j].Name())
			if extension != ".go" {
				continue
			}
			contents, err := ioutil.ReadFile(path.Join(dirsToCheck[i], files[j].Name()))
			if err != nil {
				return nil, err
			}
			match := copyrightRegex.Match(contents)
			if !match {
				badFiles = append(badFiles, path.Join(dirsToCheck[i], files[j].Name()))
			}
		}
	}
	return badFiles, nil
}
