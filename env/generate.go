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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mysteriumnetwork/go-ci/github"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
)

const ppaDevReleaseVersion = "0.0.0"

type EnvVar struct {
	Key BuildVar
	Val string
}

var buildTime = time.Now().UTC()

// GenerateEnvFile for sourcing in other stages
func GenerateEnvFile() error {
	version, err := buildVersion()
	if err != nil {
		return err
	}
	ppaVersion, err := ppaVersion()
	if err != nil {
		return err
	}
	vars := []EnvVar{
		{TagBuild, strconv.FormatBool(isTag())},
		{SnapshotBuild, strconv.FormatBool(isSnapshot())},
		{PRBuild, strconv.FormatBool(isPR())},
		{BuildVersion, version},
		{PPAVersion, ppaVersion},
		{BuildNumber, Str(BuildNumber)},
		{GithubOwner, Str(GithubOwner)},
		{GithubRepository, Str(GithubRepository)},
		{GithubSnapshotRepository, Str(GithubSnapshotRepository)},
	}
	return WriteEnvVars(vars, "./build/env.sh")
}

func isTag() bool {
	return Str(BuildTag) != ""
}

func isSnapshot() bool {
	return Str(BuildBranch) == "master" && !isTag()
}

func isPR() bool {
	return !isSnapshot() && !isTag()
}

func ppaVersion() (string, error) {
	if isTag() {
		return Str(BuildTag), EnsureEnvVars(BuildTag)
	}
	return ppaDevReleaseVersion, nil
}

func buildVersion() (string, error) {
	if isTag() {
		return Str(BuildTag), EnsureEnvVars(BuildTag)
	}
	if isPR() {
		// TODO find a format for branch version, perhaps similar to snapshot?
		return fmt.Sprintf("0.0.0-branch%s", Str(BuildBranchSafe)), nil
	}
	return snapshotVersion()
}

func snapshotVersion() (string, error) {
	if err := EnsureEnvVars(GithubOwner, GithubRepository, GithubAPIToken); err != nil {
		return "", err
	}
	releaser, err := github.NewReleaser(Str(GithubOwner), Str(GithubRepository), Str(GithubAPIToken))
	if err != nil {
		return "", err
	}
	latestRelease, err := releaser.Latest()
	if err != nil {
		return "", err
	} else if latestRelease == nil {
		return "", errors.Errorf("could not find latest release in githubRepo %s/%s", Str(GithubOwner), Str(GithubRepository))
	}
	gitLocalRepo, err := git.PlainOpen("./")
	if err != nil {
		return "", err
	}
	gitHead, err := gitLocalRepo.Head()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-1snapshot-%s-%s",
		latestRelease.TagName,
		buildTime.Format("20060102T1504"),
		gitHead.Hash().String()[:8]), nil
}

// WriteEnvVars writes vars to a shell script so they can be sourced `source env.sh` in latter build stages
func WriteEnvVars(vars []EnvVar, filename string) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, v := range vars {
		_, err := fmt.Fprintf(file, "export %s=%s;\n", v.Key, v.Val)
		if err != nil {
			return err
		}
	}
	return nil
}
