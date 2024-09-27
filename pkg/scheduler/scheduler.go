/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"os"
	"time"
	"xpanse-agent/pkg/config"
	"xpanse-agent/pkg/executor"
	"xpanse-agent/pkg/logger"
)

func StartPolling() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to create scheduler with error: %s", err.Error()))
		os.Exit(1)
	}
	_, jobErr := scheduler.NewJob(gocron.DurationJob(time.Duration(config.LoadedConfig.PollingFrequency)*time.Second),
		gocron.NewTask(executor.PollXpanseApiAndExecuteChanges),
		gocron.WithSingletonMode(gocron.LimitMode(1)),
		gocron.WithIdentifier(uuid.MustParse(config.LoadedConfig.ServiceId)),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	if jobErr != nil {
		logger.Logger.Error(fmt.Sprintf("job scheduler failed to start with error: %s", jobErr.Error()))
		os.Exit(1)
	}
	logger.Logger.Info(fmt.Sprintf("job scheduler started for serviceId %s and resourceName %s", config.LoadedConfig.ServiceId, config.LoadedConfig.ResourceName))

	scheduler.Start()

	select {}
}
