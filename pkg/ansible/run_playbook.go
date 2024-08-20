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
	"xpanse-agent/pkg/config"
	"xpanse-agent/pkg/logger"
)

func RunPlaybook(playbookName string,
	extraVars map[string]interface{},
	inventory string,
	virtualEnvRootDir string,
	pythonVersion float32,
	manageVirtualEnv bool,
	requirementsFileNameInRepo string) error {
	var res *results.AnsiblePlaybookJSONResults
	var err error
	buff := new(bytes.Buffer)
	buffError := new(bytes.Buffer)

	if manageVirtualEnv {
		logger.Logger.Info("preparing virtual environment in " + virtualEnvRootDir)
		err = CreateVirtualEnv(virtualEnvRootDir, pythonVersion, requirementsFileNameInRepo)
		if err != nil {
			return err
		}
	}

	logger.Logger.Info("Running ansible task using ansible installed in venv " + virtualEnvRootDir)

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Become:    true,
		Inventory: inventory,
		ExtraVars: extraVars,
	}

	// constructs the ansible command to be executed.
	playbookCmd := playbook.NewAnsiblePlaybookCmd(
		playbook.WithPlaybooks(playbookName),
		playbook.WithPlaybookOptions(ansiblePlaybookOptions),
		playbook.WithBinary(fmt.Sprintf("%s/bin/ansible-playbook", virtualEnvRootDir)),
	)

	// execute the ansible command constructed above.
	exec := stdoutcallback.NewJSONStdoutCallbackExecute(
		execute.NewDefaultExecute(
			execute.WithCmd(playbookCmd),
			execute.WithCmdRunDir(config.LoadedConfig.RepoCheckoutLocation),
			execute.WithErrorEnrich(playbook.NewAnsiblePlaybookErrorEnrich()),
			execute.WithWrite(io.Writer(buff)),
			execute.WithWriteError(io.Writer(buffError)),
		),
	)
	err = exec.Execute(context.TODO())
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
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
		return err
	}

	parseAndLogAnsibleOutputForResults(res)

	return err
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
