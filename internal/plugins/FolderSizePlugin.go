package plugins

import (
	"math"
	"os"
	"path/filepath"
)

type FolderSizePlugin struct {
	PathToFolder string
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

	return map[string]interface{}{
		"path": p.PathToFolder,
		"size": (math.Round(float64(totalSize)/1024/1024/1024) * 100) / 100,
	}, err
}
