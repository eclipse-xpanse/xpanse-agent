/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"os"
	"time"
	"xpanse-agent/pkg/executor"
	"xpanse-agent/pkg/logger"
)

func StartPolling(serviceId string, resourceName string, jobFrequency int) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		os.Exit(1)
	}
	_, err = scheduler.NewJob(gocron.DurationJob(time.Duration(jobFrequency)*time.Second),
		gocron.NewTask(executor.ExecuteCommand, "docker", "ps"), // to be changed to poll xpanse API.
		gocron.WithSingletonMode(gocron.LimitMode(1)))

	if err != nil {
		return
	}
	logger.Logger.Info(fmt.Sprintf("scheduler started for serviceId %s and resourceName %s", serviceId, resourceName))

	scheduler.Start()

	select {}
}
