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

package git

import (
	"fmt"
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

func (gc *GitCommiter) Checkout(branchName string) error {
	var err error
	gc.repo, err = gogit.PlainOpen("./")
	if err != nil {
		return err
	}
	fmt.Println("repo opened")
	w, err := gc.repo.Worktree()
	if err != nil {
		return err
	}
	fmt.Println("worktree fetched")

	fmt.Println("checking out master")
	err = w.Checkout(&gogit.CheckoutOptions{
		Create: false,
		Force:  true,
	})
	if err != nil {
		return err
	}
	gc.w = w
	gc.branch = branchName
	fmt.Println("master checked out")
	return nil
}

func (gc *GitCommiter) Commit(message string, files ...string) (plumbing.Hash, error) {
	fmt.Println("adding changes")
	for _, file := range files {
		fmt.Printf("adding %q\n", file)
		_, err := gc.w.Add(file)
		if err != nil {
			return [20]byte{}, err
		}
		fmt.Printf("%q added!\n", file)
	}
	fmt.Println("changes added")

	fmt.Println("performing commit")
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
	fmt.Println("Commit done")
	return commitHash, nil
}

func (gc *GitCommiter) Tag(tagVersion string, hash plumbing.Hash) error {
	fmt.Println("Tagging...", tagVersion)
	n := plumbing.ReferenceName("refs/tags/" + tagVersion)
	t := plumbing.NewHashReference(n, hash)
	err := gc.repo.Storer.SetReference(t)
	if err != nil {
		return err
	}
	fmt.Println("tagged")
	return nil
}

func (gc *GitCommiter) Push() error {
	fmt.Println("Pushing...")
	rs := config.RefSpec("refs/tags/*:refs/tags/*")
	rsm := config.RefSpec(fmt.Sprintf("refs/heads/%v:refs/heads/%v", gc.branch, gc.branch))
	err := gc.repo.Push(&gogit.PushOptions{
		RemoteName: "origin",
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
	fmt.Println("Push done")
	return nil
}
