package pluginConfigurationLoader

import (
	"encoding/json"
	"fmt"
	"os"
)

type PluginConfig struct {
	Name       string
	Parameters map[string]interface{}
}

func GetPluginsConfiguration() ([]PluginConfig, bool) {
	pluginsConfig, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return nil, true
	}
	return pluginsConfig, false
}

func loadConfig() ([]PluginConfig, error) {
	file, err := os.Open("appsettings.json")
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var data struct {
		Plugins []PluginConfig
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding config file: %v", err)
	}

	return data.Plugins, nil
}
