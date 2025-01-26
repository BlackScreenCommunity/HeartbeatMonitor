package plugins

import (
	"math"
	"os"
	"path/filepath"
)

type FolderSizePlugin struct {
	PathToFolder string
	WarningValue float64
}

func (p FolderSizePlugin) Name() string {
	return "FolderSizePlugin"
}

func (p FolderSizePlugin) Collect() (map[string]interface{}, error) {
	var totalSize int64

	err := filepath.Walk(p.PathToFolder, func(_ string, info os.FileInfo, err error) error {
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
		"path": p.PathToFolder,
		"size": (math.Round(float64(totalSize)/1024/1024/1024) * 100) / 100,
	}

	result["isWarning"] = result["size"].(float64) > p.WarningValue

	return result, err
}
