package pluginDispatcher

import (
	"fmt"
	"project/internal/pluginConfigurationLoader"
	"project/internal/pluginFactory"
	"project/internal/plugins"
)

var registeredPlugins = make(map[string]plugins.Plugin)

func InitializePlugins(pluginsConfig []pluginConfigurationLoader.PluginConfig) {
	fmt.Println("Loaded plugin configurations:")
	for _, plugin := range pluginsConfig {
		fmt.Printf("Name: %s\n",
			plugin.Name)
	}

	for _, cfg := range pluginsConfig {
		plugin, err := pluginFactory.CreatePlugin(cfg.Name, cfg.Parameters)
		if err != nil {
			fmt.Printf("Error creating plugin: %v\n", err)
			continue
		}

		RegisterPlugin(plugin)
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

func CollectAll() {
	for name, plugin := range registeredPlugins {
		data, err := plugin.Collect()
		if err != nil {
			fmt.Printf("Error collecting data from plugin %s: %v\n", name, err)
			continue
		}
		fmt.Printf("=== %s ===\n", plugin.Name())
		for key, value := range data {
			fmt.Printf("%-15s: %v\n", key, value)
		}
		fmt.Println()
	}
}
