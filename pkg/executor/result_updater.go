/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package executor

import (
	"context"
	results "github.com/apenella/go-ansible/v2/pkg/execute/result/json"
	"github.com/google/uuid"
	"xpanse-agent/pkg/logger"
	"xpanse-agent/pkg/xpanseclient"
)

func updateResultToXpanse(err error, ansibleOutput *results.AnsiblePlaybookJSONResults, changeId uuid.UUID) {
	var isSuccessful bool
	var errorString string
	var tasks []xpanseclient.AnsibleTaskResult

	if err != nil {
		isSuccessful = false
		errorString = err.Error()
	} else {
		isSuccessful = true
		errorString = ""
	}
	if ansibleOutput != nil {
		for _, play := range ansibleOutput.Plays {
			for _, task := range play.Tasks {
				name := task.Task.Name
				for _, result := range task.Hosts {
					taskSuccessful := !result.Failed
					var resultMessage string
					if result.Msg != nil {
						resultMessage = result.Msg.(string)
					} else {
						resultMessage = ""
					}
					tasks = append(tasks, xpanseclient.AnsibleTaskResult{
						IsSuccessful: taskSuccessful,
						Message:      &resultMessage,
						Name:         name,
					})
				}
			}
		}
	}

	c, clientError := getXpanseApiClient()
	if clientError != nil {
		logger.Logger.Error(clientError.Error())
	}

	if c != nil {
		resp, requestError := c.UpdateServiceChangeResultWithResponse(context.Background(), changeId, xpanseclient.ServiceChangeResult{
			Error:        &errorString,
			IsSuccessful: isSuccessful,
			Tasks:        &tasks,
		})
		if requestError != nil {
			logger.Logger.Error(requestError.Error())
		}
		if resp != nil {
			logger.Logger.Info(resp.Status())
		}
	}
}
