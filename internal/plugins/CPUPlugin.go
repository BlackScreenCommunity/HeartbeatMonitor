package plugins

type CPUPlugin struct{}

func (c CPUPlugin) Name() string {
	return "CPU Plugin"
}

func (c CPUPlugin) Collect() (map[string]interface{}, error) {
	return map[string]interface{}{
		"cpu_usage": "15%",
	}, nil
}
