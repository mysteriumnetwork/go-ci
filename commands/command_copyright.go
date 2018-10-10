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
	"regexp"

	"github.com/mysteriumnetwork/go-ci/util"
)

var copyrightRegex = regexp.MustCompile(`Copyright \(C\) \d{4} The "MysteriumNetwork/`)

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
		fmt.Println("go list crashed")
		return err
	}
	badFiles, err := getFilesWithoutCopyright(res)
	if err != nil {
		return err
	}
	if len(badFiles) != 0 {
		fmt.Println("Following files are missing copyright headers:")
		for _, v := range badFiles {
			fmt.Println(v)
		}
		return errors.New("Missing copyright headers")
	}
	fmt.Println("All files have required copyright headers!")
	return nil
}
