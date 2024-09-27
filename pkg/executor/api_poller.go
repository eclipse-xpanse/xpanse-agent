/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package executor

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"xpanse-agent/pkg/config"
	"xpanse-agent/pkg/logger"
	"xpanse-agent/pkg/xpanseclient"
)

func PollXpanseApiAndExecuteChanges() error {
	c, clientError := getXpanseApiClient()
	if clientError != nil {
		logger.Logger.Error(clientError.Error())
	}

	if c != nil {
		resp, requestError := c.GetPendingConfigurationChangeRequestWithResponse(context.Background(), config.LoadedConfig.ServiceId, config.LoadedConfig.ResourceName)
		if requestError != nil {
			logger.Logger.Error(requestError.Error())
			os.Exit(1)
		}

		if resp.StatusCode() == http.StatusNoContent {
			logger.Logger.Error("no pending configuration update requests found")
		} else if resp.StatusCode() != http.StatusOK {
			logger.Logger.Error(fmt.Sprintf("Expected HTTP 200 but received %d", resp.StatusCode()))
		}

		if resp.JSON400 != nil {
			logger.Logger.Error(strings.Join(resp.JSON400.Details, ", "))
			if resp.JSON400.ResultType == xpanseclient.ParametersInvalid {
				logger.Logger.Error("Exiting agent. Agent not started with valid parameters")
				os.Exit(1)
			}
		}
		if resp.JSON400 != nil {
			logger.Logger.Error(strings.Join(resp.JSON400.Details, ", "))
			if resp.JSON400.ResultType == xpanseclient.ParametersInvalid {
				logger.Logger.Error("Exiting agent. Agent not started with valid parameters")
				os.Exit(1)
			}
		}

		if resp.JSON200 != nil {
			logger.Logger.Info(fmt.Sprintf("pending configuration change request received for change Id %s", resp.JSON200.ChangeId))
			err := ConfigUpdate(*resp.JSON200)
			if err != nil {
				logger.Logger.Error(err.Error())
				return err
			}
		}
	}
	return nil
}
