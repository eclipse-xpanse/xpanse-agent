/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package poller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"xpanse-agent/pkg/logger"
	"xpanse-agent/pkg/xpanseclient"
)

func PollXpanseApi(serviceId string, resourceName string, xpanseApiEndpoint string) {
	hc := &http.Client{
		Transport: &logger.HttpRequestLogger{
			RoundTripper: http.DefaultTransport,
		},
	}
	{
		c, err := xpanseclient.NewClientWithResponses(xpanseApiEndpoint, xpanseclient.WithHTTPClient(hc))
		if err != nil {
			logger.Logger.Error(err.Error())
		}

		if c != nil {
			resp, requestError := c.GetPendingConfigurationChangeRequestWithResponse(context.Background(), serviceId, resourceName)
			if requestError != nil {
				logger.Logger.Error(requestError.Error())
				os.Exit(1)
			}

			if resp.StatusCode() != http.StatusOK {
				logger.Logger.Error(fmt.Sprintf("Expected HTTP 200 but received %d", resp.StatusCode()))
			}

			if resp.JSON400 != nil {
				logger.Logger.Error(strings.Join(resp.JSON400.Details, ", "))
				if resp.JSON400.ResultType == xpanseclient.ParametersInvalid {
					logger.Logger.Error("Exiting agent. Agent not started with valid parameters")
					os.Exit(1)
				}
			}
			logger.Logger.Info(fmt.Sprintf("resp.JSON200: %v", resp.JSON200))
		}

	}
}
