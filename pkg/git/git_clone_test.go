/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloneProject(t *testing.T) {
	err := CloneProject("https://github.com/swaroopar/osc", "test", "feature/testAnsible")
	if err != nil {
		assert.Empty(t, err)
	}
}

func TestCloneProjectError(t *testing.T) {
	err := CloneProject("https://github.com/swaroopar/osc-12345", "test", "feature/testAnsible")
	if err != nil {
		assert.NotEmpty(t, err)
		assert.Error(t, err, git.ErrRepositoryNotExists)
	}
}
