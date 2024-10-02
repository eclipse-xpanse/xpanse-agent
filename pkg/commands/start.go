/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"xpanse-agent/pkg/config"
	"xpanse-agent/pkg/logger"
	"xpanse-agent/pkg/scheduler"
)

var configFile string

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "To start xpanse-agent",
	Run: func(cmd *cobra.Command, args []string) {
		if config.LoadedConfig.ServiceId == "" {
			fmt.Println("ServiceId parameter is missing. Exiting.")
			os.Exit(1)
		}

		if config.LoadedConfig.ResourceName == "" {
			fmt.Println("ResourceName parameter is missing. Exiting.")
			os.Exit(1)
		}

		if config.LoadedConfig.XpanseApiEndpoint == "" {
			fmt.Println("XpanseApiEndpoint parameter is missing. Exiting.")
			os.Exit(1)
		}

		logger.Logger.Info(fmt.Sprintf("Agent started with serviceId %s, resourceName %s and with polling frequency %d",
			config.LoadedConfig.ServiceId, config.LoadedConfig.ResourceName, config.LoadedConfig.PollingFrequency))

		scheduler.StartPolling()
	},

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(Logo)
		// load config before running the command.
		errorLoading := config.LoadConfig(configFile)
		return errorLoading
	},
}

func init() {
	RootCmd.AddCommand(startCommand)
	startCommand.Flags().StringP(serviceIdArgument, "s", "", "ID of the service. Mandatory parameter.")
	startCommand.Flags().StringP(resourceNameArgument, "r", "", "ID of the resource on which the agent is running. Mandatory parameter.")
	startCommand.Flags().StringP(xpanseApiEndpointArgument, "e", "", "Xpanse API end point to poll for configuration changes. Mandatory parameter.")
	startCommand.Flags().IntP(jobFrequencyArgument, "p", config.DefaultApiPollingFrequencyInSeconds, "Polling frequency for change requests in seconds")
	startCommand.Flags().StringVarP(&configFile, configFileLocationArgument, "c", config.DefaultAgentConfigFile, "Config file can contain all input arguments or some. "+
		"This parameter must contain the relative path of the config file. Eg: /tmp/config.yaml")
	startCommand.Flags().StringP(repoCheckoutFolderArgument, "f", config.DefaultRepoCheckoutFolder, "Folder to which GIT repos will be checked out to.")

	// necessary to bind the flag input of cobra command to viper config.
	// then all data can be accessed via the AgentConfig struct.
	bindInputArgumentsToViper()
}

func bindInputArgumentsToViper() {
	_ = viper.BindPFlag(config.ServiceIdKey, startCommand.Flags().Lookup(serviceIdArgument))
	_ = viper.BindPFlag(config.ResourceNameKey, startCommand.Flags().Lookup(resourceNameArgument))
	_ = viper.BindPFlag(config.PollingFrequencyKey, startCommand.Flags().Lookup(jobFrequencyArgument))
	_ = viper.BindPFlag(config.XpanseApiEndpointKey, startCommand.Flags().Lookup(xpanseApiEndpointArgument))
	_ = viper.BindPFlag(config.RepoCheckoutLocation, startCommand.Flags().Lookup(repoCheckoutFolderArgument))
}
