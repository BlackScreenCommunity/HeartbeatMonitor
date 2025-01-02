package main

import (
	"fmt"
	"project/internal/plugins"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type PluginConfig struct {
	Name string `json:"name"`
}

var registeredPlugins = make(map[string]plugins.Plugin)

func main() {

	pluginsConfig, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}

	fmt.Println("Loaded plugin configurations:")
	for _, plugin := range pluginsConfig {
		fmt.Printf("Name: %s\n",
			plugin.Name)
	}

	for _, cfg := range pluginsConfig {
		plugin, err := plugins.CreatePlugin(cfg.Name)
		if err != nil {
			fmt.Printf("Error creating plugin: %v\n", err)
			continue
		}

		RegisterPlugin(plugin)
	}

	CollectAll()
}

func loadConfig() ([]PluginConfig, error) {
	viper.SetConfigName("appsettings")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	var plugins []PluginConfig
	if err := viper.UnmarshalKey("plugins", &plugins); err != nil {
		return nil, fmt.Errorf("error unmarshalling plugins: %v", err)
	}

	return plugins, nil
}

func RegisterPlugins(plugins []plugins.Plugin) {
	for _, registrationPlugin := range plugins {
		RegisterPlugin(registrationPlugin)
	}
}

func RegisterPlugin(p plugins.Plugin) {
	registeredPlugins[uuid.New().String()] = p
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
