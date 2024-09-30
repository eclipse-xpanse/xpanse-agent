/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	agentConfig "xpanse-agent/pkg/config"
)

func init() {
	agentConfig.LoadedConfig = &agentConfig.AgentConfig{
		ServiceId:            "test",
		ResourceName:         "test",
		RepoCheckoutLocation: "/tmp",
		PollingFrequency:     0,
		XpanseApiEndpoint:    "",
	}

}
func TestCloneProject(t *testing.T) {
	t.Setenv("XPANSE_AGENT_SERVICE_ID", "test")
	err := CloneRepository("https://github.com/swaroopar/osc", "feature/testAnsible")
	if err != nil {
		assert.Empty(t, err)
	}

}

func TestCloneProjectError(t *testing.T) {
	err := CloneRepository("https://github.com/swaroopar/osc-12345", "feature/testAnsible")
	if err != nil {
		assert.NotEmpty(t, err)
		assert.Error(t, err, git.ErrRepositoryNotExists)
	}
}
