/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

func CloneProject(projectUrl string, serviceId string, branch string) error {
	gitConfig := &git.CloneOptions{
		URL:      projectUrl,
		Progress: os.Stdout,
	}
	var repository *git.Repository
	var err error

	repository, err = git.PlainClone(fmt.Sprintf("/tmp/%s", serviceId), false, gitConfig)

	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) && repository == nil {
			repository, err = git.PlainOpen(fmt.Sprintf("/tmp/%s", serviceId))
			if err != nil {
				fmt.Println("Error opening repository:", err)
				return err
			}
		}
	}

	if repository != nil {
		// Fetch updates from the remote
		err = repository.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{"refs/*:refs/*"},
		})
		if err != nil {
			if !errors.Is(err, git.NoErrAlreadyUpToDate) {
				fmt.Println("Error fetching updates:", err)
				return err
			}
		}

		// load default main/master worktree
		w, workTreeError := repository.Worktree()

		checkoutOptions := &git.CheckoutOptions{Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)), Force: true}
		if w != nil {
			checkoutErr := w.Checkout(checkoutOptions)
			if checkoutErr != nil {
				return checkoutErr
			}
		}
		if workTreeError != nil {
			return workTreeError
		}
	}

	return nil
}
