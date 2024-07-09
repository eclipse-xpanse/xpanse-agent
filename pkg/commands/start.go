/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"xpanse-agent/pkg/scheduler"
)

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "To start xpanse-agent",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag(ServiceIdArgument).Value.String() == "" {
			fmt.Println("ServiceId is missing")
			os.Exit(1)
		}

		frequency, errNew := strconv.Atoi(cmd.Flag(JobFrequencyArgument).Value.String())
		if errNew != nil {
			return
		}
		scheduler.StartPolling(
			cmd.Flag(CommandArgument).Value.String(),
			frequency,
			cmd.Flag(CommandArgumentsArgument).Value.String(),
			cmd.Flag(ServiceIdArgument).Value.String())
	},

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(Logo)
		return initializeConfig(cmd)
	},
}

func init() {
	RootCmd.AddCommand(startCommand)
	startCommand.Flags().StringP(ServiceIdArgument, "s", "", "ID of the service. Mandatory parameter.")
	startCommand.Flags().StringP(CommandArgument, "c", "ansible-pull", "command that must be executed.")
	startCommand.Flags().StringP(CommandArgumentsArgument, "a", "", "command that must be executed.")
	startCommand.Flags().IntP(JobFrequencyArgument, "j", 30, "Job Frequency in seconds")
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetEnvPrefix("XPANSE_AGENT")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				return
			}
		}
	})
}
