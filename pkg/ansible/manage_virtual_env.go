/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package ansible

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"xpanse-agent/pkg/git"
	"xpanse-agent/pkg/logger"
)

func createVirtualEnv(virtualEnvDir string, pythonVersion string, moduleRequirementsFile string) error {
	var cmd *exec.Cmd
	var err error
	// Check if the virtualenv already exists
	if _, err = os.Stat(virtualEnvDir); !os.IsNotExist(err) {
		logger.Logger.Info(fmt.Sprintf("Virtual environment '%s' already exists. Installing required modules on the same.", virtualEnvDir))
	} else {
		cmd = exec.Command(fmt.Sprintf("python%s", pythonVersion), "-m", "venv", virtualEnvDir)
		err = cmd.Run()
	}
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error creating virtual environment: %s", err))
		return err
	}

	if moduleRequirementsFile != "" {
		err = installRequirements(virtualEnvDir, moduleRequirementsFile)
	}
	if err != nil {
		logger.Logger.Info(fmt.Sprintf("Virtual environment '%s' prepared successfully.", virtualEnvDir))
	}
	return err
}

func installRequirements(virtualEnvDir string, requirementsFile string) error {
	logger.Logger.Info(fmt.Sprintf("Installing required modules in %s", virtualEnvDir))
	var err error
	originalWorkingDir, err := os.Getwd()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error getting current working directory: %s", err))
		return err
	}
	err = os.Chdir(git.GetRepoDirectory())
	logger.Logger.Info(fmt.Sprintf("running pip install from %s", git.GetRepoDirectory()))
	if err != nil {
		return err
	}
	if _, err = os.Stat(requirementsFile); os.IsNotExist(err) {
		logger.Logger.Error(fmt.Sprintf("Requirements file '%s' does not exist.", requirementsFile))
		return err
	}

	// Prepare the pip install command
	cmd := exec.Command(fmt.Sprintf("%s/bin/pip", virtualEnvDir), "install", "-r", requirementsFile)

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error installing packages: %s", err))
		return err
	}

	logger.Logger.Info("Required modules installed successfully")
	logger.Logger.Info(string(output))

	// return back to original directory.
	err = os.Chdir(originalWorkingDir)
	return err
}

func GetVirtualEnvRootDirectory(virtualEnvRootDir string) string {
	var usedVirtualEnvVar string
	if virtualEnvRootDir != "" {
		if strings.HasSuffix(virtualEnvRootDir, "/") {
			usedVirtualEnvVar = strings.TrimSuffix(virtualEnvRootDir, "/")
		} else {
			usedVirtualEnvVar = virtualEnvRootDir
		}
	} else {
		usedVirtualEnvVar = "/tmp/virtualEnv"
	}
	return usedVirtualEnvVar
}
