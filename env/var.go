// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package env

// BuildVar env variables used in CI build.
// Some of them will be calculated when generating env file, others should be passed though env.
type BuildVar string

const (
	// TagBuild indicates release build
	TagBuild = BuildVar("RELEASE_BUILD")

	// RCBuild indicates a release candidate build (containing "-rc")
	RCBuild = BuildVar("RC_BUILD")

	// SnapshotBuild indicates snapshot release build (master branch)
	SnapshotBuild = BuildVar("SNAPSHOT_BUILD")

	// PRBuild indicates pull-request build
	PRBuild = BuildVar("PR_BUILD")

	// BuildVersion stores build version
	BuildVersion = BuildVar("BUILD_VERSION")

	// PPAVersion stores build version for PPA
	PPAVersion = BuildVar("PPA_VERSION")

	// BuildNumber stores CI build number
	BuildNumber = BuildVar("BUILD_NUMBER")

	// BuildTag stores git tag for build
	BuildTag = BuildVar("BUILD_TAG")

	// BuildBranch stores branch name
	BuildBranch = BuildVar("BUILD_BRANCH")

	// BuildBranchSafe stores branch name, with special characters replaced with hyphens
	BuildBranchSafe = BuildVar("BUILD_BRANCH_SAFE")

	// GithubOwner stores github repository's owner
	GithubOwner = BuildVar("GITHUB_OWNER")

	// GithubRepository stores github repository name
	GithubRepository = BuildVar("GITHUB_REPO")

	// GithubSnapshotRepository stores github repository name for snapshot builds
	GithubSnapshotRepository = BuildVar("GITHUB_SNAPSHOT_REPO")

	// GithubAPIToken is used for accessing github API
	GithubAPIToken = BuildVar("GITHUB_API_TOKEN")

	// DockerHubUsername is hub.docker.com username under which to push snapshot builds
	DockerHubUsername = BuildVar("DOCKERHUB_USERNAME")

	// DockerHubPassword is hub.docker.com password of DockerHubUsername
	DockerHubPassword = BuildVar("DOCKERHUB_PASSWORD")

	// CommitMessage is HEAD commit message
	CommitMessage = BuildVar("CI_COMMIT_MESSAGE")
)
