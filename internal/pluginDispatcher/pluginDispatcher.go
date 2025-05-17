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

// Creates exemplar of plugins and
// adds to a list of registered plugins
func InitializePlugins(pluginsConfig []PluginConfig) {
	fmt.Println("Loaded plugin configurations:")

	for _, pluginConfig := range pluginsConfig {
		if pluginConfig.Active {
			fmt.Printf("Create plugin: %s\n", pluginConfig.Name)
			plugin, err := pluginFactory.CreatePlugin(pluginConfig.Name, pluginConfig.Parameters)
			if err != nil {
				fmt.Printf("Error creating plugin: %v\n", err)
				continue
			}
			RegisterPlugin(plugin)
		}
	}
}

// Adds current plugin's exemplar
// into a list of registered plugins
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

// Returns list of registered and active plugins
func GetPlugins() map[string]plugins.Plugin {
	return registeredPlugins
}
