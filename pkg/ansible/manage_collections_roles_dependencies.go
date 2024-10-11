/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package ansible

import (
	"context"
	"fmt"
	"github.com/apenella/go-ansible/v2/pkg/execute"
	galaxycollectioninstall "github.com/apenella/go-ansible/v2/pkg/galaxy/collection/install"
	galaxyroleinstall "github.com/apenella/go-ansible/v2/pkg/galaxy/role/install"
	"xpanse-agent/pkg/git"
	"xpanse-agent/pkg/logger"
)

func installGalaxyDependencies(requirementsFile string, virtualEnv string) error {
	// install collections from requirements file
	galaxyInstallCollectionCmd := galaxycollectioninstall.NewAnsibleGalaxyCollectionInstallCmd(
		galaxycollectioninstall.WithBinary(fmt.Sprintf("%s/bin/ansible-galaxy",
			GetVirtualEnvRootDirectory(virtualEnv))),
		galaxycollectioninstall.WithGalaxyCollectionInstallOptions(
			&galaxycollectioninstall.AnsibleGalaxyCollectionInstallOptions{
				ForceWithDeps:    true,
				RequirementsFile: git.GetRepoDirectory() + "/" + requirementsFile,
			}))
	galaxyInstallCollectionExec := execute.NewDefaultExecute(
		execute.WithCmd(galaxyInstallCollectionCmd),
	)
	err := galaxyInstallCollectionExec.Execute(context.TODO())
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}

	// install roles from requirements file
	galaxyInstallRoleCmd := galaxyroleinstall.NewAnsibleGalaxyRoleInstallCmd(
		galaxyroleinstall.WithBinary(fmt.Sprintf("%s/bin/ansible-galaxy",
			GetVirtualEnvRootDirectory(virtualEnv))),
		galaxyroleinstall.WithGalaxyRoleInstallOptions(
			&galaxyroleinstall.AnsibleGalaxyRoleInstallOptions{
				ForceWithDeps: true,
				RoleFile:      git.GetRepoDirectory() + "/" + requirementsFile,
			}))
	galaxyInstallRoleExec := execute.NewDefaultExecute(
		execute.WithCmd(galaxyInstallRoleCmd),
	)
	err = galaxyInstallRoleExec.Execute(context.TODO())
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}
	return nil
}
