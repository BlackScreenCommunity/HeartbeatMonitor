package plugins

import "time"

type ServerTimePlugin struct{}

func (v ServerTimePlugin) Name() string {
	return "ServerTimePlugin"
}

func (p ServerTimePlugin) Collect() (map[string]interface{}, error) {
	return map[string]interface{}{
		"time": time.Now().Format("02-01-2006 15:04:05"),
	}, nil
}
