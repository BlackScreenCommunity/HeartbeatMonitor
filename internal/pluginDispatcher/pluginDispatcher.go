package pluginDispatcher

import (
	"encoding/json"
	"fmt"

	"project/internal/pluginFactory"
	"project/internal/plugins"
)

type PluginConfig struct {
	Name       string
	Active     bool
	Parameters map[string]interface{}
}

// Collection of registered plugins for these agent
var registeredPlugins = make(map[string]plugins.Plugin)

func InitializePlugins(pluginsConfig []PluginConfig) {
	fmt.Println("Loaded plugin configurations:")
	for _, plugin := range pluginsConfig {
		fmt.Printf("Name: %s\n",
			plugin.Name)
	}

	for _, cfg := range pluginsConfig {
		if cfg.Active {
			plugin, err := pluginFactory.CreatePlugin(cfg.Name, cfg.Parameters)
			if err != nil {
				fmt.Printf("Error creating plugin: %v\n", err)
				continue
			}
			RegisterPlugin(plugin)
		}
	}
}

func RegisterPlugins(plugins []plugins.Plugin) {
	for _, registrationPlugin := range plugins {
		RegisterPlugin(registrationPlugin)
	}
}

func RegisterPlugin(p plugins.Plugin) {
	pluginNumber := fmt.Sprintf("%04d", len(registeredPlugins)+1)

	registeredPlugins[pluginNumber+"_"+p.Name()] = p
}

// Gets plugins metrics
func CollectAll() map[string]interface{} {

	pluginResultCollection := make(map[string]interface{})

	for name, plugin := range registeredPlugins {
		data, err := plugin.Collect()
		if err != nil {
			fmt.Printf("Error collecting data from plugin %s: %v\n", name, err)
			continue
		}

		pluginResultCollection[plugin.Name()] = data
	}
	return pluginResultCollection
}

// Returns metrics from all plugins in JSON format
func GetPluginsJsonData() string {
	jsonData, err := json.MarshalIndent(CollectAll(), "", "  ")

	if err != nil {
		return ""
	}

	return string(jsonData)
}

func GetPlugins() map[string]plugins.Plugin {
	return registeredPlugins
}
