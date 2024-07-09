/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"xpanse-agent/version"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "To view the version of the xpanse-agent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("xpanse-agent version : %s \n", version.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCommand)
}
