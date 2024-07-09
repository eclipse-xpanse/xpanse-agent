/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package executor

import (
	"bufio"
	"io"
	"os/exec"
	"xpanse-agent/pkg/logger"
)

func ExecuteCommand(command string, arguments string) {
	cmd := exec.Command(command, arguments)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	go func() {
		merged := io.MultiReader(stdout, stderr)
		scanner := bufio.NewScanner(merged)
		for scanner.Scan() {
			logger.Logger.Info(scanner.Text())
		}
	}()

	startError := cmd.Start()
	if startError != nil {
		logger.Logger.Error(startError.Error())
	}

	waitError := cmd.Wait()

	if waitError != nil {
		logger.Logger.Error(waitError.Error())
	}

}
