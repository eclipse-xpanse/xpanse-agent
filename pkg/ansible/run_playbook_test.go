/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package ansible

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"xpanse-agent/pkg/config"
)

func TestRunPlaybook(t *testing.T) {
	err := config.LoadConfig("test_data/test-xpanse-agent-config.yml")
	assert.Nil(t, err)
	results, runError := RunPlaybook(
		"kafka-container-manage.yml", nil, "", "/tmp/kafka-test/", 3.10, true, "requirements.txt")
	assert.Nil(t, runError)
	assert.NotNil(t, results)
}
