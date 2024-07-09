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

func StartPolling(command string, jobFrequency int, arguments string, serviceId string) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		os.Exit(1)
	}
	_, err = scheduler.NewJob(gocron.DurationJob(time.Duration(jobFrequency)*time.Second),
		gocron.NewTask(executor.ExecuteCommand, command, arguments),
		gocron.WithSingletonMode(gocron.LimitMode(1)))

	if err != nil {
		return
	}
	logger.Logger.Info(fmt.Sprintf("xpanse-agent started for serviceId %s", serviceId))

	scheduler.Start()

	select {}
}
