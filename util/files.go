// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

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
