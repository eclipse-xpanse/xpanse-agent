/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package executor

import (
	"fmt"
	results "github.com/apenella/go-ansible/v2/pkg/execute/result/json"
	"xpanse-agent/pkg/ansible"
	"xpanse-agent/pkg/git"
	"xpanse-agent/pkg/logger"
	"xpanse-agent/pkg/xpanseclient"
)

func informPanicErrorToXpanse(changeId string) {
	if r := recover(); r != nil {
		panicError := fmt.Errorf("change request failed with error: %s", r)
		logger.Logger.Error(panicError.Error())
		updateResultToXpanse(panicError, nil, changeId)
	}
}

func ConfigUpdate(request xpanseclient.ServiceConfigurationChangeRequest) error {
	defer informPanicErrorToXpanse(request.ChangeId.String())
	var err error
	var result *results.AnsiblePlaybookJSONResults
	err = git.CloneProject(request.AnsibleScriptConfig.RepoUrl,
		request.AnsibleScriptConfig.Branch)
	if err == nil {
		result, err = ansible.RunPlaybook(request.AnsibleScriptConfig.PlaybookName,
			*request.ConfigParameters,
			"",
			request.AnsibleScriptConfig.VirtualEnv,
			request.AnsibleScriptConfig.PythonVersion,
			true,
			request.AnsibleScriptConfig.RequirementsFile,
		)
	}
	updateResultToXpanse(err, result, request.ChangeId.String())
	return err
}
