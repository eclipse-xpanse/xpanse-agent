/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDefaultValues(t *testing.T) {
	errEnv := os.Setenv("XPANSE_AGENT_SERVICE_ID", "test_value")
	assert.Nil(t, errEnv)
	err := LoadConfig("test_data/not-existing.yml")
	assert.Nil(t, err)
	assert.Equal(t, "/tmp", LoadedConfig.RepoCheckoutLocation)
	assert.Equal(t, 30, LoadedConfig.PollingFrequency)
	assert.Equal(t, "test_value", LoadedConfig.ServiceId)
	assert.Equal(t, "", LoadedConfig.XpanseApiEndpoint)
}

func TestConfigFromFile(t *testing.T) {
	errEnv := os.Setenv("XPANSE_AGENT_SERVICE_ID", "test_value")
	assert.Nil(t, errEnv)
	err := LoadConfig("test_data/test-xpanse-agent-config.yml")
	assert.Nil(t, err)
	assert.Equal(t, "/tmp", LoadedConfig.RepoCheckoutLocation)
	assert.Equal(t, 40, LoadedConfig.PollingFrequency)
	assert.Equal(t, "test_value", LoadedConfig.ServiceId)
	assert.Equal(t, "", LoadedConfig.ResourceName)
}
