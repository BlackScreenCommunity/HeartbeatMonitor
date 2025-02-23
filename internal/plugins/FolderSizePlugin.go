package plugins

import (
	"math"
	"os"
	"path/filepath"
)

// FolderSizePlugin is a plugin for calculating the total size of files in the specified directory.
// - PathToFolder: path to the directory to scan.
// - WarningValue: threshold value in GB; if exceeded, the isWarning flag is set.
// - InstanceName: name of the current plugin instance.
type FolderSizePlugin struct {
	PathToFolder string
	WarningValue float64
	InstanceName string
}

// Name returns the name of the current plugin instance.
func (pluginConfig FolderSizePlugin) Name() string {
	return pluginConfig.InstanceName
}

// Calculates the total size of the folder
// Set the 'IsWarning' flag when the total size is bigger than the WarningValue
func (pluginConfig FolderSizePlugin) Collect() (map[string]interface{}, error) {
	var totalSize int64

	err := filepath.Walk(pluginConfig.PathToFolder, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return map[string]interface{}{
			"isWarning": true,
			"Message":   "Path not found",
		}, nil
	}

	result := map[string]interface{}{
		"path": pluginConfig.PathToFolder,
		"size": (math.Round(float64(totalSize)/1024/1024/1024) * 100) / 100,
	}

	result["isWarning"] = result["size"].(float64) > pluginConfig.WarningValue

	return result, err
}
