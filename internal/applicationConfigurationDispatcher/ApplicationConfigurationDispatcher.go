package applicationConfigurationDispatcher

import (
	"encoding/json"
	"fmt"
	"os"
	"project/internal/pluginDispatcher"
)

// Defines the local agent name
type ServerInfo struct {
	Name      string
	Login     string
	Password  string
	AgentMode bool
}

// Parameters for the web interface
// and network interaction with external agents
type WebServerConfig struct {
	Port                int
	AgentPollingTimeout int
}

// Remote agent parameters
type AgentConfig struct {
	Address  string
	Name     string
	Active   bool
	Login    string
	Password string
}

// Overall set of application parameters
type ApplicationConfiguration struct {
	Server    ServerInfo
	WebServer WebServerConfig
	Plugins   []pluginDispatcher.PluginConfig
	Agents    []AgentConfig
}

// Reads application's configuration
// from appsettings.json
func GetConfigFromFile() (ApplicationConfiguration, bool) {
	pluginsConfig, err := loadConfigFromFile()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return ApplicationConfiguration{}, true
	}
	return pluginsConfig, false
}

// Handles operations with file.
// Tries to search and decode JSON object to
// Application Configuration
func loadConfigFromFile() (ApplicationConfiguration, error) {
	file, err := os.Open("appsettings.json")
	if err != nil {
		return ApplicationConfiguration{}, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var data ApplicationConfiguration
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return ApplicationConfiguration{}, fmt.Errorf("error decoding config file: %v", err)
	}

	return data, nil
}
