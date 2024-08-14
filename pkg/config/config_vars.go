/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package config

const (
	XpanseEnvironmentVarPrefix = "XPANSE_AGENT"

	// any changes to key names must be manually updated in the AgentConfig mapstructure.

	ServiceIdKey         = "service_id"
	ResourceNameKey      = "resource_name"
	PollingFrequencyKey  = "polling_frequency"
	XpanseApiEndpointKey = "xpanse_api_endpoint"
	RepoCheckoutLocation = "repo_checkout_location"

	DefaultRepoCheckoutFolder           = "/tmp"
	DefaultApiPollingFrequencyInSeconds = 30
	DefaultAgentConfigFile              = "./xpanse-agent-config.yml"
)
