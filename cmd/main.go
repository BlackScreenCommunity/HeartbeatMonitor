package main // определили пакет main

import (
    "fmt"
    "project/internal/plugins"
)

var registeredPlugins = make(map[string]plugins.Plugin)

func main() {
	RegisterPlugin(plugins.VersionPlugin{})
	RegisterPlugin(plugins.CPUPlugin{})
	CollectAll()
}

func RegisterPlugin(p plugins.Plugin) {
    registeredPlugins[p.Name()] = p
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