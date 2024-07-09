/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package commands

import (
	"github.com/spf13/cobra"
	"os"
)

var Logo = `
                                                              _
__  ___ __   __ _ _ __  ___  ___        __ _  __ _  ___ _ __ | |_
\ \/ / '_ \ / _' | '_ \/ __|/ _ \_____ / _' |/ _' |/ _ \ '_ \| __|
 >  <| |_) | (_| | | | \__ \  __/_____| (_| | (_| |  __/ | | | |_
/_/\_\ .__/ \__,_|_| |_|___/\___|      \__,_|\__, |\___|_| |_|\__|
     |_|                                     |___/
`

var RootCmd = &cobra.Command{
	Use:  "xpanse-agent",
	Long: Logo,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				return
			}
			os.Exit(0)
		}
	},
}
