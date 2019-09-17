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

package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsPathExcluded determines if the provided path is excluded from common searches
func IsPathExcluded(paths []string, path string) bool {
	for _, exclude := range paths {
		if strings.HasPrefix(path, exclude) {
			return true
		}
	}
	return false
}

// GetPackagePathsWithExcludes returns the package paths gotten by go list minus the dirs excluded
func GetPackagePathsWithExcludes(path string, excludes ...string) ([]string, error) {
	res, err := GetPackagePaths(path)
	if err != nil {
		return res, err
	}

	result := make([]string, 0)
	for _, dir := range res {
		isExcluded := false
		for _, v := range excludes {
			if strings.Contains(dir, "/"+v) {
				isExcluded = true
			}
		}
		if !isExcluded {
			result = append(result, dir)
		}
	}

	return result, nil
}

// GetPackagePaths gets the go paths for various checks
func GetPackagePaths(path string) ([]string, error) {
	cmd := exec.Command("go", "list", path)
	res, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}
	stringified := string(res)
	splits := strings.Split(stringified, "\n")
	return splits, nil
}

// GoLintExcludes returns commonly excluded dirs from quality checks
func GoLintExcludes() []string {
	return []string{
		".idea",
		".git",
		"build",
		"vendor",
	}
}

// GetProjectFileDirectories returns all the project directories excluding git and vendor
func GetProjectFileDirectories(paths []string) ([]string, error) {
	directories := make([]string, 0)

	root := "./"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !IsPathExcluded(paths, path) && path != root {
			directories = append(directories, path)
		}
		return nil
	})
	return directories, err
}
