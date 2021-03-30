// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package github

import (
	"context"
	"log"
	"os"
	"path/filepath"

	gogithub "github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

// Releaser releases to Github
type Releaser struct {
	client            *gogithub.Client
	owner, repository string
}

// Release represents a Github release
type Release struct {
	ID      int64
	TagName string
	*Releaser
}

// NewReleaser creates a new Releaser instance
func NewReleaser(owner, repo, token string) (*Releaser, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(context.Background(), ts)
	return &Releaser{
		client:     gogithub.NewClient(oauthClient),
		owner:      owner,
		repository: repo,
	}, nil
}

// Create creates a new Github release
func (r *Releaser) Create(name string) (*Release, error) {
	release, _, err := r.client.Repositories.CreateRelease(context.Background(), r.owner, r.repository, &gogithub.RepositoryRelease{Name: gogithub.String(name), TagName: gogithub.String(name)})
	if err != nil {
		return nil, err
	}
	log.Printf("created release ID: %v, tag: %v", *release.ID, *release.TagName)
	return &Release{ID: *release.ID, TagName: *release.TagName, Releaser: r}, nil
}

// Find finds Github release by tag name
func (r *Releaser) Find(tagName string) (*Release, error) {
	release, _, err := r.client.Repositories.GetReleaseByTag(context.Background(), r.owner, r.repository, tagName)
	if err != nil {
		return nil, err
	}
	log.Printf("found release ID: %v, tag: %v", *release.ID, *release.TagName)
	return &Release{ID: *release.ID, TagName: *release.TagName, Releaser: r}, nil
}

// Latest finds the latest Github release
func (r *Releaser) Latest() (*Release, error) {
	release, _, err := r.client.Repositories.GetLatestRelease(context.Background(), r.owner, r.repository)
	if err != nil {
		return nil, err
	}
	return &Release{ID: *release.ID, TagName: *release.TagName, Releaser: r}, nil
}

// UploadAsset uploads asset to the release
func (r *Release) UploadAsset(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	asset, _, err := r.client.Repositories.UploadReleaseAsset(context.Background(), r.owner, r.repository, r.ID, &gogithub.UploadOptions{
		Name: filepath.Base(file.Name()),
	}, file)
	if err != nil {
		return err
	}
	log.Println("uploaded asset ", *asset.Name)
	return nil
}
