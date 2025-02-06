package plugins

import (
	"math"
	"os"
	"path/filepath"
)

type FolderSizePlugin struct {
	PathToFolder string
	WarningValue float64
	InstanceName string
}

func (ppluginConfig FolderSizePlugin) Name() string {
	return ppluginConfig.InstanceName
}

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
