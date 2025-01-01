package main

import (
    "fmt"
    "github.com/google/uuid"
    "project/internal/plugins"
)

var registeredPlugins = make(map[string]plugins.Plugin)

func main() {
	pluginCollection := []plugins.Plugin{
        plugins.VersionPlugin{},
        plugins.CPUPlugin{},
		plugins.HardDriveFreeSpacePlugin{},

        plugins.PostgreSqlQueryPlugin{
            ConnectionString: "host=localhost port=5432 user=postgres password=bpmsoft dbname=01-pulse sslmode=disable",
            Query:            "SELECT Count(\"Id\") FROM \"Contact\"",
            PluginName:       "Contacts count",
        },

        plugins.PostgreSqlQueryPlugin{
            ConnectionString: "host=localhost port=5432 user=postgres password=bpmsoft dbname=01-pulse sslmode=disable",
            Query:            "SELECT Count(\"Id\") FROM \"SysPrcPersistentStore\"",
            PluginName:       "Process logs count",
        },
    }

	RegisterPlugins(pluginCollection)
	CollectAll()
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
