// Copyright (c) 2020 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package git

import (
	"fmt"
	"log"
	"time"

	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type GitCommiter struct {
	repo   *gogit.Repository
	w      *gogit.Worktree
	branch string
	token  string
}

func NewCommiter(apiToken string) *GitCommiter {
	return &GitCommiter{
		token: apiToken,
	}
}

type CheckoutOptions struct {
	BranchName string
	Force      bool
	Keep       bool
}

func (gc *GitCommiter) Checkout(options *CheckoutOptions) error {
	var err error
	gc.repo, err = gogit.PlainOpen("./")
	if err != nil {
		return err
	}
	log.Println("repo opened")
	w, err := gc.repo.Worktree()
	if err != nil {
		return err
	}
	log.Println("worktree fetched")

	log.Println("checking out master")
	err = w.Checkout(&gogit.CheckoutOptions{
		Create: false,
		Force:  options.Force,
		Keep:   options.Keep,
	})
	if err != nil {
		return err
	}
	gc.w = w
	gc.branch = options.BranchName
	log.Println("master checked out")
	return nil
}

func (gc *GitCommiter) Commit(message string, files ...string) (plumbing.Hash, error) {
	log.Println("adding changes")
	for _, file := range files {
		log.Printf("adding %q\n", file)
		_, err := gc.w.Add(file)
		if err != nil {
			return [20]byte{}, err
		}
		log.Printf("%q added!\n", file)
	}
	log.Println("changes added")

	log.Println("performing commit")
	commitHash, err := gc.w.Commit(message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "Mister CI tool",
			Email: "dev@mysterium.network",
			When:  time.Now(),
		},
	})
	if err != nil {
		return commitHash, err
	}
	log.Println("Commit done")
	return commitHash, nil
}

func (gc *GitCommiter) Tag(tagVersion string, hash plumbing.Hash) error {
	log.Println("Tagging...", tagVersion)
	n := plumbing.ReferenceName("refs/tags/" + tagVersion)
	t := plumbing.NewHashReference(n, hash)
	err := gc.repo.Storer.SetReference(t)
	if err != nil {
		return err
	}
	log.Println("tagged")
	return nil
}

type PushOptions struct {
	Remote string
}

func (gc *GitCommiter) Push(options *PushOptions) error {
	log.Println("Pushing...")
	rs := config.RefSpec("refs/tags/*:refs/tags/*")
	rsm := config.RefSpec(fmt.Sprintf("refs/heads/%v:refs/heads/%v", gc.branch, gc.branch))
	err := gc.repo.Push(&gogit.PushOptions{
		RemoteName: options.Remote,
		Auth: &http.BasicAuth{
			// this can be anything but not an empty string
			Username: "MisterFancyPants",
			Password: gc.token,
		},
		RefSpecs: []config.RefSpec{rs, rsm},
	})
	if err != nil {
		return err
	}
	log.Println("Push done")
	return nil
}
