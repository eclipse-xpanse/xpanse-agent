/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package ansible

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apenella/go-ansible/v2/pkg/execute"
	results "github.com/apenella/go-ansible/v2/pkg/execute/result/json"
	"github.com/apenella/go-ansible/v2/pkg/execute/stdoutcallback"
	"github.com/apenella/go-ansible/v2/pkg/playbook"
	"io"
	"strings"
	"xpanse-agent/pkg/git"
	"xpanse-agent/pkg/logger"
)

func RunPlaybook(playbookName string,
	extraVars map[string]interface{},
	inventory *map[string]interface{},
	virtualEnvRootDir string,
	pythonVersion float32,
	manageVirtualEnv bool,
	requirementsFileNameInRepo string) (*results.AnsiblePlaybookJSONResults, error) {
	var res *results.AnsiblePlaybookJSONResults
	var err error
	buff := new(bytes.Buffer)
	buffError := new(bytes.Buffer)
	var usedVirtualEnvVar string
	var inventoryFileName string
	if virtualEnvRootDir != "" {
		if strings.HasSuffix(virtualEnvRootDir, "/") {
			usedVirtualEnvVar = strings.TrimSuffix(virtualEnvRootDir, "/")
		} else {
			usedVirtualEnvVar = virtualEnvRootDir
		}
	} else {
		usedVirtualEnvVar = "/tmp/virtualEnv"
	}
	if manageVirtualEnv {
		logger.Logger.Info("preparing virtual environment in " + usedVirtualEnvVar)
		err = createVirtualEnv(usedVirtualEnvVar, pythonVersion, requirementsFileNameInRepo)
		if err != nil {
			return nil, err
		}
	}

	logger.Logger.Info("Running ansible task using ansible installed in venv " + usedVirtualEnvVar)
	if inventory != nil {
		inventoryFile, err := createInventoryFile(inventory)
		if err != nil {
			return nil, err
		}
		inventoryFileName = inventoryFile.Name()
	} else {
		inventoryFileName = ""
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Become:    true,
		Inventory: inventoryFileName,
		ExtraVars: extraVars,
	}

	// constructs the ansible command to be executed.
	playbookCmd := playbook.NewAnsiblePlaybookCmd(
		playbook.WithPlaybooks(playbookName),
		playbook.WithPlaybookOptions(ansiblePlaybookOptions),
		playbook.WithBinary(fmt.Sprintf("%s/bin/ansible-playbook", usedVirtualEnvVar)),
	)

	// execute the ansible command constructed above.
	exec := stdoutcallback.NewJSONStdoutCallbackExecute(
		execute.NewDefaultExecute(
			execute.WithCmd(playbookCmd),
			execute.WithCmdRunDir(git.GetRepoDirectory()),
			execute.WithErrorEnrich(playbook.NewAnsiblePlaybookErrorEnrich()),
			execute.WithWrite(io.Writer(buff)),
			execute.WithWriteError(io.Writer(buffError)),
		),
	)
	err = exec.Execute(context.TODO())
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	// all warnings from Ansible are written to stderr stream.
	if buffError.Len() > 0 {
		for _, line := range bytes.Split(buffError.Bytes(), []byte("\n")) {
			if len(line) > 0 {
				logger.Logger.Warn(string(line))
			}
		}
	}

	res, err = results.ParseJSONResultsStream(io.Reader(buff))
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	parseAndLogAnsibleOutputForResults(res)
	deleteInventoryFile(inventoryFileName)
	return res, err
}

func parseAndLogAnsibleOutputForResults(ansibleOutput *results.AnsiblePlaybookJSONResults) {
	for _, play := range ansibleOutput.Plays {
		for _, task := range play.Tasks {
			name := task.Task.Name
			for host, result := range task.Hosts {
				logger.Logger.Info("|" + host + "|" + play.Play.Name + "|" + name + "|" + fmt.Sprintf("%v", result.StdoutLines) + "|" + fmt.Sprintf("%t", result.Changed) + "|" + fmt.Sprintf("%t", result.Failed))
			}
		}
	}
}
