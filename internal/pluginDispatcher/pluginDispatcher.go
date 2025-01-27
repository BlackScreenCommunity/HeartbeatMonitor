package pluginDispatcher

import (
	"encoding/json"
	"fmt"

	"project/internal/pluginFactory"
	"project/internal/plugins"
	"project/internal/utils"
)

type PluginConfig struct {
	Name       string
	Active     bool
	Parameters map[string]interface{}
}

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

func CollectAll() map[string]interface{} {

	pluginResultCollection := make(map[string]interface{})

	for name, plugin := range registeredPlugins {
		data, err := plugin.Collect()
		if err != nil {
			fmt.Printf("Error collecting data from plugin %s: %v\n", name, err)
			continue
		}
		if len(data) == 1 {
			for _, value := range data {
				pluginResultCollection[plugin.Name()] = value
			}
		} else {
			pluginResultCollection[plugin.Name()] = data
		}
	}
	return utils.MapDereference(pluginResultCollection)
}

func GetPluginsJsonData() string {
	jsonData, err := json.MarshalIndent(CollectAll(), "", "  ")

	if err != nil {
		return ""
	}

	return string(jsonData)
}

// TODO Поправить вывод в консоль
func PrintPluginResult(pluginResultCollection map[string]interface{}) {
	for pluginResult := range pluginResultCollection {
		// fmt.Printf("=== %s ===\n", plugin.Name())
		for key, value := range pluginResult {
			fmt.Printf("%-15s: %v\n", key, value)
		}
		fmt.Println()
	}
}
