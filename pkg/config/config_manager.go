/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"xpanse-agent/pkg/logger"
)

type AgentConfig struct {
	ServiceId            string `mapstructure:"service_id"`
	ResourceName         string `mapstructure:"resource_name"`
	RepoCheckoutLocation string `mapstructure:"repo_checkout_location"`
	PollingFrequency     int    `mapstructure:"polling_frequency"`
	XpanseApiEndpoint    string `mapstructure:"xpanse_api_endpoint"`
}

var LoadedConfig AgentConfig

func LoadConfig(configFileLocation string) error {
	setDefaultValues()
	loadConfigFromEnvironmentVariables()
	loadConfigurationFromFile(configFileLocation)

	// load config from all sources into AgentConfig struct.
	if err := viper.Unmarshal(&LoadedConfig); err != nil {
		return err
	}
	return nil
}

func splitConfig(config string) (string, string, string) {
	dir, file := filepath.Split(config)
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	filename := strings.TrimSuffix(file, filepath.Ext(file))
	return dir, filename, ext
}

func setDefaultValues() {
	viper.SetDefault(ServiceIdKey, "")
	viper.SetDefault(ResourceNameKey, "")
	viper.SetDefault(PollingFrequencyKey, DefaultApiPollingFrequencyInSeconds)
	viper.SetDefault(RepoCheckoutLocation, DefaultRepoCheckoutFolder)
	viper.SetDefault(XpanseApiEndpointKey, "")
}

func loadConfigFromEnvironmentVariables() {
	viper.SetEnvPrefix(XpanseEnvironmentVarPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func loadConfigurationFromFile(configFileLocation string) {
	logger.Logger.Info("Loading agent configuration from " + configFileLocation)
	dir, filename, extension := splitConfig(configFileLocation)
	viper.AddConfigPath(dir)
	viper.SetConfigName(filename)
	viper.SetConfigType(extension)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
		// Config file not found, use default values
		fmt.Println("No config file provided or not found, using default values")
	}
}
