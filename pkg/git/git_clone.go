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
	agentConfig "xpanse-agent/pkg/config"
	"xpanse-agent/pkg/logger"
)

func CloneProject(projectUrl string, branch string) error {
	gitConfig := &git.CloneOptions{
		URL:      projectUrl,
		Progress: nil,
	}
	var repository *git.Repository
	var err error
	var w *git.Worktree

	repository, err = git.PlainClone(GetRepoDirectory(), false, gitConfig)

	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) && repository == nil {
			logger.Logger.Info("Repository already cloned. Reusing it.")
			repository, err = git.PlainOpen(GetRepoDirectory())
			if err != nil {
				fmt.Println("Error opening repository:", err)
				return err
			}
		}
	}

	if repository != nil {
		// Fetch updates from the remote
		logger.Logger.Info("Fetching updates from remote.")
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
		w, err = repository.Worktree()

		checkoutOptions := &git.CheckoutOptions{Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)), Force: true}
		if w != nil {
			checkoutErr := w.Checkout(checkoutOptions)
			if checkoutErr != nil {
				return checkoutErr
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func GetRepoDirectory() string {
	return agentConfig.LoadedConfig.RepoCheckoutLocation + "/" + agentConfig.LoadedConfig.ServiceId
}
