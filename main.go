/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package main

import (
	"os"
	"xpanse-agent/pkg/commands"
)

func main() {
	err := commands.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
