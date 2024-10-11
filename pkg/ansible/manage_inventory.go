/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package ansible

import (
	"encoding/json"
	"fmt"
	"os"
	"xpanse-agent/pkg/config"
	"xpanse-agent/pkg/logger"
)

func createInventoryFile(inventoryContent *map[string]interface{}) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "*-"+config.LoadedConfig.ServiceId+"-inventory.json")
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error creating file: %s", err))
		return nil, err
	}

	jsonData, err := json.Marshal(&inventoryContent)
	if err != nil {
		return nil, err
	}

	_, err = tempFile.Write(jsonData)
	if err != nil {
		return nil, err
	}

	logger.Logger.Info(fmt.Sprintf("Inventory file created at %s", tempFile.Name()))
	return tempFile, nil
}

func deleteInventoryFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		logger.Logger.Error("Error deleting inventory file. ignoring the error.")
	}
	logger.Logger.Info("temporary inventory file cleaned up.")
}
