package plugins

import (
	"fmt"
	"math"

	"github.com/shirou/gopsutil/v3/disk"
)

type HardDriveFreeSpacePlugin struct {
	DriveMountPoint       string
	WarningValueInGb      float64
	WarningValueInPercent float64
}

func (v HardDriveFreeSpacePlugin) Name() string {
	return "HardDriveFreeSpacePlugin"
}

func (plgn HardDriveFreeSpacePlugin) Collect() (map[string]interface{}, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions: %v", err)
	}

	result := make(map[string]interface{})
	for _, partition := range partitions {
		if plgn.DriveMountPoint != partition.Mountpoint {
			continue
		}
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to get usage for partition %s: %v", partition.Mountpoint, err)
		}

		partitionTotals := map[string]interface{}{
			"total":       (math.Round(float64(usage.Total)/1024/1024/1024) * 100) / 100, // Общее пространство (в байтах)
			"free":        (math.Round(float64(usage.Free)/1024/1024/1024) * 100) / 100,  // Свободное пространство (в байтах)
			"used":        (math.Round(float64(usage.Used)/1024/1024/1024) * 100) / 100,  // Используемое пространство (в байтах)
			"usedPercent": (math.Round(usage.UsedPercent) * 100) / 100,                   // Процент использования
		}

		partitionTotals["isWarning"] = (partitionTotals["free"].(float64) < plgn.WarningValueInGb) || (partitionTotals["usedPercent"].(float64) < plgn.WarningValueInPercent)
		result[partition.Mountpoint] = partitionTotals
	}

	isWarning := false

	for _, value := range result {

		isWarning = isWarning || value.(map[string]interface{})["isWarning"].(bool)
	}

	result["isWarning"] = isWarning
	return result, nil
}
