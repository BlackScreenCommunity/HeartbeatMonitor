package main

import (
	"encoding/json"
	"fmt"
	"os"
	"project/internal/pluginFactory"
	"project/internal/plugins"

	"github.com/google/uuid"
)

type PluginConfig struct {
	Name       string
	Parameters map[string]interface{}
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
		plugin, err := pluginFactory.CreatePlugin(cfg.Name, cfg.Parameters)
		if err != nil {
			fmt.Printf("Error creating plugin: %v\n", err)
			continue
		}

		RegisterPlugin(plugin)
	}

	CollectAll()
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
