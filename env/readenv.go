// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package env

import (
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

// Bool reads a bool env var.
// EnsureEnvVars should be called first to ensure it has a specified value.
func Bool(v BuildVar) bool {
	env := os.Getenv(string(v))
	val, _ := strconv.ParseBool(env)
	return val
}

// Str reads a string env var.
// EnsureEnvVars should be called first to ensure it has a specified value.
func Str(v BuildVar) string {
	return os.Getenv(string(v))
}

// IfRelease performs func passed as an arg if current build is any kind of release
func IfRelease(do func() error) error {
	isRelease, err := IsRelease()
	if err != nil {
		return err
	}
	if isRelease {
		log.Println("release build detected, performing conditional action")
		return do()
	}
	log.Println("not a release build, skipping conditional action")
	return nil
}

// IsRelease true when building a release (tag or snapshot).
func IsRelease() (bool, error) {
	if err := EnsureEnvVars(TagBuild, SnapshotBuild); err != nil {
		return false, err
	}
	return Bool(TagBuild) || Bool(SnapshotBuild), nil
}

// IsPR true when building a Pull Request.
func IsPR() (bool, error) {
	if err := EnsureEnvVars(PRBuild); err != nil {
		return false, err
	}
	return Bool(PRBuild), nil
}

// IsFullBuild true when full build is requested via commit message `[ci full]`
func IsFullBuild() (bool, error) {
	if err := EnsureEnvVars(CommitMessage); err != nil {
		return false, err
	}
	match, err := regexp.MatchString(`\[ci full]`, Str(CommitMessage))
	return match, errors.Wrap(err, "failed to parse commit message")
}
