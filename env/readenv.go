/*
 * Copyright (C) 2019 The "MysteriumNetwork/go-ci" Authors.
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
