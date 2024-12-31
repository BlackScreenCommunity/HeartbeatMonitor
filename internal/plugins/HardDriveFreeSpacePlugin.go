package plugins

import (
    "fmt"
    "math"
    "strings"
    "github.com/shirou/gopsutil/v3/disk"
)

type HardDriveFreeSpacePlugin struct{}

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
         usage, err := disk.Usage(partition.Mountpoint)
         if err != nil {
             return nil, fmt.Errorf("failed to get usage for partition %s: %v", partition.Mountpoint, err)
         }
         
         partitionTotals := map[string]interface{}{
            "total":       (math.Floor(float64(usage.Total) / 1024 / 1024 / 1024) * 100) /100 ,       // Общее пространство (в байтах)
            "free":        (math.Floor(float64(usage.Free)  / 1024 / 1024 / 1024) * 100) /100 ,        // Свободное пространство (в байтах)
            "used":        (math.Floor(float64(usage.Used)  / 1024 / 1024 / 1024) * 100) /100 ,        // Используемое пространство (в байтах)
            "usedPercent": (math.Floor(usage.UsedPercent) * 100) /100 , // Процент использования
        }
        
        var sb strings.Builder
        
        for key, value := range partitionTotals {
			sb.WriteString(fmt.Sprintf("%s: %v\t", key, value))
		 }

         result[partition.Mountpoint] = sb.String()
     }
 
     return result, nil
}