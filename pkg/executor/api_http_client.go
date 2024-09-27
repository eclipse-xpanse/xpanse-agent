/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package executor

import (
	"net/http"
	"xpanse-agent/pkg/config"
	"xpanse-agent/pkg/logger"
	"xpanse-agent/pkg/xpanseclient"
)

var xpanseHttpClient *xpanseclient.ClientWithResponses

// singleton client generator
func getXpanseApiClient() (*xpanseclient.ClientWithResponses, error) {
	var err error
	if xpanseHttpClient == nil {
		hc := &http.Client{
			Transport: &logger.HttpRequestLogger{
				RoundTripper: http.DefaultTransport,
			},
		}
		xpanseHttpClient, err = xpanseclient.NewClientWithResponses(config.LoadedConfig.XpanseApiEndpoint, xpanseclient.WithHTTPClient(hc))
	}
	return xpanseHttpClient, err
}
